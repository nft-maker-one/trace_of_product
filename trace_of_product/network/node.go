package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/database"
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var NodeTables = []Node{}

const BlockInterval = 3

// time.Second 只可以与常数连乘

type Node struct {
	Id   int
	Addr string
}

type NodeServer struct {
	Node
	Db                  *database.NodeDb
	Interval            int
	Chain               *core.Blockchain
	MessagePool         map[types.Hash]core.Block
	PrePareConfirmCount map[types.Hash]map[int]bool
	CommitConfirmCount  map[types.Hash]map[int]bool
	IsCommitBroadcast   map[types.Hash]bool
	isReply             map[types.Hash]bool
	priKey              *crypto.PrivateKey
	PubKey              crypto.PublicKey
	SequenceId          int
	Pool                *MemoryPool
	lock                sync.Mutex
	IsLeader            bool
}

func NewNodeServer(addr string) *NodeServer {
	id := randomId()
	server := &NodeServer{}
	server.Id = id
	cfg, err := utils.NewConfig("./config.yaml")
	if err != nil {
		logrus.Errorln("read config error")
		panic(err)
	}
	server.Db = database.InitNodeDb(cfg.Mysql.Dsn)
	server.Addr = addr
	server.PubKey = crypto.GenerateKeyPair(id)
	server.MessagePool = make(map[types.Hash]core.Block)
	server.PrePareConfirmCount = make(map[types.Hash]map[int]bool)
	server.CommitConfirmCount = make(map[types.Hash]map[int]bool)
	server.IsCommitBroadcast = make(map[types.Hash]bool)
	server.isReply = make(map[types.Hash]bool)
	server.priKey, err = crypto.ReadPriKey(id)
	server.Interval = BlockInterval
	server.Pool = NewMemoryPool(500)
	if err != nil {
		panic(err)
	}
	leaderId, err := server.Db.AddNode(server.Id, server.Addr, server.PubKey)
	if err != nil {
		panic(err)
	}
	NodeTables = append(NodeTables, Node{
		Id:   server.Id,
		Addr: server.Addr,
	})
	if leaderId == server.Id {
		server.IsLeader = true
	}
	server.Chain = core.CreateChain(leaderId)
	server.lock = sync.Mutex{}

	utils.LogMsg([]string{"NewNodeServer"}, []string{"create server successfully"})
	return server

}

func (s *NodeServer) handleRequest(data []byte) {
	rpc := &RPC{}
	err := json.Unmarshal(data, rpc)
	if err != nil {
		utils.LogError([]string{"handleRequest"}, []string{err.Error()})
	}
	switch rpc.ContentType {
	case cRequest:
		s.handleClientRequest(rpc.Payload) //payload = message
	case cPrePrepare:
		s.handlePrePrepare(rpc.Payload) //payload = preprepare
	case cPrepare:
		s.handlePrepare(rpc.Payload) //payload = prepare
	case cCommit:
		s.handleCommit(rpc.Payload)
	case cEgg:
		s.handleEgg(rpc.Payload)
	case cTest:
		message := Message{}
		err := json.Unmarshal(rpc.Payload, &message)
		if err != nil {
			fmt.Printf("%s receive message successfully message:%s", s.Addr, string(data))
		}
	}

}

func (s *NodeServer) CreateBlock() {
	time.Sleep(BlockInterval * 60 * time.Second)
	eggs := []*core.Eggplant{}
	if s.Pool.fifo.len > s.Pool.Cap {
		for i := 0; i < s.Pool.Cap; i++ {
			egg, _ := s.Pool.PopEgg()
			eggs = append(eggs, &egg)
		}
	} else {
		for {
			egg, err := s.Pool.PopEgg()
			if err != nil {
				break
			}
			eggs = append(eggs, &egg)
		}
	}
	dataHash, err := core.CalculateDataHash(eggs)
	if err != nil {
		utils.LogMsg([]string{"CreateBlock"}, []string{"data hash error " + err.Error()})
		return
	}
	h := &core.Header{
		Version:       0,
		PrevBlockHash: s.Chain.GetPrevHash(),
		DataHash:      dataHash,
		Timestamp:     time.Now().Unix(),
		Height:        int32(s.Chain.Height()) + 1,
		Nonce:         rand.Int63n(math.MaxInt64),
	}
	h.UpdateScore(s.Chain, eggs)
	h.SelectLeader()
	if h.Leader == 0 {
		user := database.ConsortiumNode{}
		s.Db.DB.Where("id != ?", s.Id).Limit(1).Find(&user)
		h.Leader = user.Id
	}
	b, err := core.NewBlock(h, eggs)
	if err != nil {
		utils.LogMsg([]string{"CreateBlock"}, []string{"create new block failed err = " + err.Error()})
		return
	}
	b.GetHash()
	err = b.Sign(*s.priKey)
	if err != nil {
		utils.LogMsg([]string{"CreateBlock"}, []string{"sign block failed err = " + err.Error()})
		return
	}
	s.SequenceId++
	utils.LogMsg([]string{"CreateBlock"}, []string{fmt.Sprintf("create new visual graph with sequenceId %d", s.SequenceId)})

	sig, err := s.priKey.Sign(b.Header.Bytes())

	pp := PrePrepare{}
	pp.RequestMessage = *(b)
	pp.SequencId = s.SequenceId
	pp.Digest = b.GetHash()
	// 1 传输前后 RequestMessage 的哈希是一致的
	if b.GetHash() != pp.RequestMessage.GetHash() {
		fmt.Printf("%+v\n", *(b))
		fmt.Printf("%+v\n", *(b.Header))
		fmt.Printf("%+v\n", pp.RequestMessage)
		fmt.Printf("%+v\n", *(pp.RequestMessage.Header))
	}
	if err != nil {
		utils.LogError([]string{"handleClientRequest"}, []string{fmt.Sprintf("sign error %s ", err.Error())})
		return
	}
	pp.Sign = sig.ToByte()
	s.MessagePool[b.Hash(core.BlockHasher{})] = *b
	ppDate, err := json.Marshal(&pp)
	if err != nil {
		utils.LogMsg([]string{"handleClientRequest"}, []string{"prePrepare marshal failed"})
		return
	}
	rpc := RPC{}
	rpc.ContentType = cPrePrepare
	rpc.Payload = ppDate
	s.Broadcast(rpc)

}

// 对 rpc 中的 Message decode，验证其中的区块信息，并对 Message 进行签名
func (s *NodeServer) handleClientRequest(payload []byte) {
	req := ClientRequest{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		utils.LogMsg([]string{"handleClientRequest"}, []string{"client request decode failed err = " + err.Error()})
	}
	switch req.Header {
	case "Upload":
		egg := new(core.Eggplant)
		egg.Decode(core.NewGobEggplantDecoder(bytes.NewReader(req.Content)))
		if egg.EggplantId <= 0 {
			utils.LogError([]string{"handleClientRequest"}, []string{"eggplant has a invalid id"})
			return
		}
		hash := egg.SetHash(core.EggplantHasher{})
		if hash != egg.Hash {
			utils.LogError([]string{"handleClientRequest"}, []string{"eggplant has a invalid hash"})
			return
		}
		err = egg.Sign(*s.priKey)
		if err != nil {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"sign error " + err.Error()})
			return
		}
		egg.NodeId = s.Id
		egg.FirstSeen = time.Now().Unix()
		buf := &bytes.Buffer{}
		err = egg.Encode(core.NewGobEggplantEncoder(buf))
		if err != nil {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"encode eggplant error " + err.Error()})
			return
		}
		s.Pool.AddEgg(*egg)
		eggMes := &EggMes{}
		eggMes.Egg = buf.Bytes()
		eggMes.NodeId = s.Id
		eggMesBuf := &bytes.Buffer{}
		err = eggMes.Encode(NewGobEggMesEncoder(eggMesBuf))
		if err != nil {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"encode eggMes error " + err.Error()})
			return
		}
		rpc := RPC{}
		rpc.ContentType = cEgg
		rpc.Payload = eggMesBuf.Bytes()
		s.Broadcast(rpc)
	case "Search":
		id, err := strconv.Atoi(string(req.Content))
		if err != nil {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"id format error " + err.Error()})
			return
		}
		for _, block := range s.Chain.Chains {
			for _, egg := range block.Eggplants {
				if egg.EggplantId == id {
					reqMeta, err := json.Marshal(&egg.MetaData)
					if err != nil {
						utils.LogMsg([]string{"handleClientRequest"}, []string{"encode matedata failed err = " + err.Error()})
						return
					}
					http.Post(req.RespAddr, "application/json", bytes.NewReader(reqMeta))
					return
				}
			}
		}
	default:
		fmt.Println("Invalid request")
	}

}

func (s *NodeServer) handlePrePrepare(payload []byte) {
	pp := &PrePrepare{}
	err := json.Unmarshal(payload, pp)
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{err.Error()})
		return
	}
	if s.SequenceId+1 != pp.SequencId {
		utils.LogError([]string{"handlePrePrepare"}, []string{fmt.Sprintf("node expected SequencdId [%d] but got [%d]", s.SequenceId+1, pp.SequencId)})
		return
	}
	hash := pp.RequestMessage.Hash(core.BlockHasher{})
	if hash != pp.Digest {
		utils.LogError([]string{"handlePrePrepare"}, []string{fmt.Sprintf("digest is not correct,expect %v, but got %v", hash, pp.Digest)})
		return
	}
	// 获取 leader 节点的公钥
	leaderId := s.Chain.Chains[s.Chain.Height()].Leader
	node, err := s.Db.SearchNodeById(leaderId)
	if err != nil {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			utils.LogMsg([]string{"handlePrePrepare"}, []string{fmt.Sprintf("SearchNodeById for %d times", i)})
			node, err = s.Db.SearchNodeById(leaderId)
			if err == nil {
				break
			}
			if i == 4 {
				utils.LogError([]string{"handlePrePrepare"}, []string{fmt.Sprintf("node [%d] can not connect to database, please send after that", node.Id)})
				return
			}
		}
	}
	// 验证是否得到公钥节点签名
	sig, err := crypto.ByteToSignature(pp.Sign)
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{"the signature has a fault in decode format"})
		return
	}
	if !sig.Verify(node.PubKey, pp.RequestMessage.Header.Bytes()) {
		utils.LogError([]string{"handlePrePrepare"}, []string{"refuse prepare,the PrePrepare Message is not signed by the Leader Node"})
		return
	}
	p := Prepare{}
	// 更新视图
	s.SequenceId++
	// 更新内存池
	s.MessagePool[pp.Digest] = pp.RequestMessage

	p.Digest = pp.Digest
	p.NodeId = s.Id
	p.SequencId = s.SequenceId
	sig, err = s.priKey.Sign(hash[:])
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{"refuse prepare,the PrePrepare Message is not signed by the Leader Node"})
		return
	}
	p.Sign = sig.ToByte()
	pDate, err := json.Marshal(&p)
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{"prepare marshal error"})
		return
	}
	rpc := RPC{}
	rpc.ContentType = cPrepare
	rpc.Payload = pDate
	s.Broadcast(rpc)

}
func (s *NodeServer) handlePrepare(payload []byte) {
	p := Prepare{}
	err := json.Unmarshal(payload, &p)
	if err != nil {
		utils.LogError([]string{"handlePrepare"}, []string{"message is not the format of prepare"})
		return
	}
	if _, ok := s.MessagePool[p.Digest]; !ok {
		utils.LogError([]string{"handlePrepare"}, []string{"the prepare message have not been stored in message pool"})
		return
	}
	if p.SequencId != s.SequenceId {
		utils.LogError([]string{"handlePrepare"}, []string{"the sequenceId is not correct"})
	}
	node, err := s.Db.SearchNodeById(p.NodeId)
	if err != nil {
		utils.LogError([]string{"handlePrepare"}, []string{"can not connect to manage system"})
		return
	}
	sign, err := crypto.ByteToSignature(p.Sign)
	if err != nil {
		utils.LogError([]string{"handlePrepare"}, []string{"the signature can not be decoded in prepare"})
		return
	}
	if !sign.Verify(node.PubKey, p.Digest[:]) {
		utils.LogError([]string{"handlePrepare"}, []string{"the signature is invalid in prepare"})
		return
	}
	s.lock.Lock()
	s.prePrepareAdd(p.Digest, p.NodeId)
	pNodeNum := len(s.PrePareConfirmCount[p.Digest])
	tNodeNum := s.Db.GetNum()
	superNodeId := s.Chain.Chains[s.Chain.Height()]
	threshold := 0
	if s.Node.Id == superNodeId.Leader {
		threshold = tNodeNum / 3 * 2
	} else {
		threshold = tNodeNum/3*2 - 1
	}
	s.lock.Unlock()
	if pNodeNum >= threshold && !s.IsCommitBroadcast[p.Digest] {
		nonce := uint64(randomId())
		signCommit := CommitData(p.Digest, nonce)
		fmt.Printf("commitData = %s", hex.EncodeToString(signCommit))
		sig, err := s.priKey.Sign(signCommit)
		if err != nil {
			utils.LogError([]string{"handlePrepare"}, []string{"sign commit data failed, err=", err.Error()})
		}
		c := new(Commit)
		c.Sign = sig.ToByte()
		c.Digest = p.Digest
		c.NodeId = s.Id
		c.SequenceId = p.SequencId
		c.Nonce = nonce
		rpc := RPC{}
		rpc.ContentType = cCommit
		payload, err = json.Marshal(c)
		if err != nil {
			utils.LogError([]string{"handlePrepare"}, []string{"marshal commit failed err= " + err.Error()})
		}
		rpc.Payload = payload
		s.Broadcast(rpc)
		s.IsCommitBroadcast[p.Digest] = true
		utils.LogMsg([]string{"handlePrepare"}, []string{"broadcast completed"})
	}

}

func (s *NodeServer) handleCommit(payload []byte) {
	c := new(Commit)
	err := json.Unmarshal(payload, &c)
	if err != nil {
		utils.LogError([]string{"handleCommit"}, []string{"payload is not the format of commit err=" + err.Error()})
	}
	if _, ok := s.PrePareConfirmCount[c.Digest]; !ok {
		utils.LogError([]string{"handleCommtit"}, []string{"the prepare is not storage in prepare pool"})
	}
	if s.SequenceId != c.SequenceId {
		utils.LogError([]string{"handleCommit"}, []string{"not the correct visual graph"})
	}
	cNode, err := s.Db.SearchNodeById(c.NodeId)
	if err != nil {
		utils.LogError([]string{"handleCommit"}, []string{"can not find the node in ConsortiumNode err=" + err.Error()})
		return
	}
	cData := CommitData(c.Digest, c.Nonce)
	sig, err := crypto.ByteToSignature(c.Sign)

	if !sig.Verify(cNode.PubKey, cData) {
		utils.LogError([]string{"handleCommit"}, []string{"signature is invalid"})
		return
	}
	s.commitAdd(c.Digest, c.NodeId)
	cNodeNum := len(s.CommitConfirmCount[c.Digest])
	tNodeNum := s.Db.GetNum()
	if cNodeNum > tNodeNum/3*2 && !s.isReply[c.Digest] && s.IsCommitBroadcast[c.Digest] {
		block := s.MessagePool[c.Digest]
		if err != nil {
			utils.LogError([]string{"handleCommit"}, []string{"decode block failed err=" + err.Error()})
			return
		}
		err = s.Chain.AddBlock(&block)
		if err != nil {
			utils.LogMsg([]string{"handleCommit"}, []string{"block is invalid err= " + err.Error()})
			return
		}
		for _, egg := range block.Eggplants {
			s.Pool.DeleteEggByHash(egg.Hash)
		}
		utils.LogMsg([]string{"handleCommit"}, []string{""})
		utils.LogMsg([]string{"handleCommit"}, []string{"add block successfully height =" + strconv.Itoa(int(block.Height)) + " hash = " + block.DataHash.String()})
		if s.Chain.GetLeader() == s.Id {
			rec := make(map[int]int)
			for _, egg := range block.Eggplants {
				rec[egg.NodeId]++
			}
			for k, v := range rec {
				err := s.Db.AddScore(k, v)
				if err != nil {
					utils.LogMsg([]string{"handleCommit"}, []string{"update verify time failed err = %v" + err.Error()})
				}
			}
			fmt.Println(rec)
			go s.CreateBlock()
		}

	}

}

func (s *NodeServer) handleEgg(payload []byte) {
	eggMes := new(EggMes)
	err := eggMes.Decode(NewGobEggMesDecoder(bytes.NewReader(payload)))
	if err != nil {
		utils.LogMsg([]string{"handleEgg"}, []string{"payload is not the format of eggMes err = " + err.Error()})
	}
	egg := new(core.Eggplant)
	err = egg.Decode(core.NewGobEggplantDecoder(bytes.NewReader(eggMes.Egg)))
	if err != nil {
		utils.LogMsg([]string{"handleEgg"}, []string{"egg decode failed, err = " + err.Error()})
		return
	}
	if egg.EggplantId <= 0 {
		utils.LogError([]string{"handleEgg"}, []string{"eggplant has a invalid id"})
		return
	}
	hash := egg.SetHash(core.EggplantHasher{})
	if hash != egg.Hash {
		utils.LogError([]string{"handleEgg"}, []string{"eggplant has a invalid hash"})
		return
	}
	if err = egg.Verify(); err != nil {
		utils.LogError([]string{"handleEgg"}, []string{"eggplant verify failed err = " + err.Error()})
		return
	}
	node, err := s.Db.SearchNodeById(eggMes.NodeId)
	if err != nil {
		utils.LogError([]string{"handleEgg"}, []string{"nodeId incorrect"})
		return
	}
	if !utils.CompareBytes(node.PubKey, egg.PublickKey) {
		utils.LogMsg([]string{"handleEgg"}, []string{"the id do not correspond with the pubkey"})
		return
	}
	s.Pool.AddEgg(*egg)
}

func (s *NodeServer) prePrepareAdd(hash types.Hash, id int) {
	if _, ok := s.PrePareConfirmCount[hash]; !ok {
		s.PrePareConfirmCount[hash] = make(map[int]bool)
	}
	s.PrePareConfirmCount[hash][id] = true
}

func (s *NodeServer) commitAdd(hash types.Hash, id int) {
	if _, ok := s.CommitConfirmCount[hash]; !ok {
		s.CommitConfirmCount[hash] = make(map[int]bool)
	}
	s.CommitConfirmCount[hash][id] = true
}

func (s *NodeServer) Broadcast(rpc RPC) {
	data, err := json.Marshal(&rpc)
	if err != nil {
		utils.LogError([]string{"broadcast"}, []string{"rpc marshal error"})
	}
	s.UpdateNodeTable()
	for _, node := range NodeTables {
		if node.Id == s.Node.Id {
			continue
		}
		tcpDial(data, node.Addr)
	}
	utils.LogMsg([]string{"broadcast"}, []string{"broadcast completed"})
}

func (s *NodeServer) SendMessageToNode(rpc RPC, id int) error {
	data, err := json.Marshal(&rpc)
	if err != nil {
		return fmt.Errorf("rpc marshal failed err = %s", err.Error())
	}
	for _, node := range NodeTables {
		if node.Id == id {
			return tcpDial(data, node.Addr)
		}
	}
	return fmt.Errorf("the node[" + strconv.Itoa(id) + "] is not storaged")
}

func (s *NodeServer) SendMessage(rpc RPC, addr string) error {
	data, err := json.Marshal(&rpc)
	if err != nil {
		return fmt.Errorf("rpc marshal failed err = " + err.Error())
	}
	return tcpDial(data, addr)
}

func (s *NodeServer) UpdateNodeTable() {
	nodes := s.Db.GetAll()
	NodeTables = []Node{}
	for _, node := range nodes {
		NodeTables = append(NodeTables, Node{Id: node.Id, Addr: node.Addr})
	}
}

func randomId() int {
	num := 0
	// 获取一个 10 位数的 id
	for num < 1000000000 {
		num = rand.Intn(10000000000)
	}
	return num
}
