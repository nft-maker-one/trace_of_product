package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/database"
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

var NodeTables = []Node{}

type Node struct {
	Id   int
	Addr string
}

type NodeServer struct {
	Node
	Db                  *database.NodeDb
	Chain               *core.Blockchain
	MessagePool         map[types.Hash]Message
	PrePareConfirmCount map[types.Hash]map[string]int
	CommitConfirmCount  map[types.Hash]map[string]int
	IsCommitBroadcast   map[types.Hash]bool
	isReply             map[types.Hash]bool
	priKey              *crypto.PrivateKey
	PubKey              crypto.PublicKey
	SequenceId          int
}

func NewNodeServer(addr string) *NodeServer {
	id := randomId()
	server := &NodeServer{}
	server.Id = id
	cfg, err := utils.NewConfig("../config.yaml")
	if err != nil {
		logrus.Errorln("read config error")
		panic(err)
	}
	server.Db = database.InitNodeDb(cfg.Mysql.Dsn)
	server.Addr = addr
	server.PubKey = crypto.GenerateKeyPair(id)
	server.MessagePool = make(map[types.Hash]Message)
	server.PrePareConfirmCount = make(map[types.Hash]map[string]int)
	server.CommitConfirmCount = make(map[types.Hash]map[string]int)
	server.isReply = make(map[types.Hash]bool)
	server.priKey, err = crypto.ReadPriKey(id)
	if err != nil {
		panic(err)
	}
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
		s.handleClientRequest(rpc.Payload)
	case cPrePrepare:
		s.handlePrePrepare(rpc.Payload)
	case cPrepare:
		s.handlePrepare(rpc.Payload)
	case cTest:
		message := Message{}
		err := json.Unmarshal(rpc.Payload, &message)
		if err != nil {
			fmt.Printf("%s receive message successfully message:%s", s.Addr, string(data))
		}
		fmt.Println("receive message from ")
	}

}

// 对 rpc 中的 Message decode，验证其中的区块信息，并对 Message 进行签名
func (s *NodeServer) handleClientRequest(payload []byte) {
	message := &Message{}
	gob.NewDecoder(bytes.NewReader(payload)).Decode(message)
	switch message.Header {
	case MessageTypeBlock:
		block := new(core.Block)
		block.Decode(core.NewGobBlockDecode(bytes.NewReader(message.Data)))
		if block.Height != int32(s.Chain.Height()) {
			utils.LogMsg([]string{"handleClientRequest"}, []string{fmt.Sprintf("expect block with height %d but got %d", s.Chain.Height(), block.Height)})
			return
		}
		hash := block.Hash(core.BlockHasher{})
		if hash != block.DataHash {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"hash not correct in the block"})
			return
		}
		if !block.Signature.Verify(block.Validator, block.DataHash[:]) {
			utils.LogMsg([]string{"handleClientRequest"}, []string{fmt.Sprintf("block [%v] has a wrong validator", block.DataHash)})
		}
		utils.LogMsg([]string{"handleClientRequest"}, []string{fmt.Sprintf("verification completed, block [%v] is valid", block.DataHash)})
		s.SequenceId++
		utils.LogMsg([]string{"handleClientRequest"}, []string{fmt.Sprintf("create new visual graph with sequenceId %d", s.SequenceId)})
		s.MessagePool[block.DataHash] = *message
		pp := PrePrepare{}
		pp.RequestMessage = *message
		pp.Digest = types.Hash(sha256.Sum256(payload))
		pp.SequencId = s.SequenceId
		sig, err := s.priKey.Sign(payload)
		if err != nil {
			utils.LogError([]string{"handleClientRequest"}, []string{fmt.Sprintf("sign error %s ", err.Error())})
		}
		pp.Sign = sig.ToByte()
		ppDate, err := json.Marshal(&pp)
		if err != nil {
			utils.LogMsg([]string{"handleClientRequest"}, []string{"prePrepare marshal failed"})
		}
		rpc := RPC{}
		rpc.ContentType = cPrePrepare
		rpc.Payload = ppDate
		s.Broadcast(rpc)
		rpc.ContentType = cPrePrepare

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
	hash := pp.RequestMessage.Hash()
	if hash != pp.Digest {
		utils.LogError([]string{"handlePrePrepare"}, []string{"digest is not correct"})
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
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), node.PubKey)
	key := crypto.PublicKey{}
	key.Key.Curve = elliptic.P256()
	key.Key.X = x
	key.Key.Y = y
	sig, err := crypto.ByteToSignature(pp.Sign)
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{"the signature has a fault in decode format"})
		return
	}
	if !sig.Verify(key, pp.Sign) {
		utils.LogError([]string{"handlePrePrepare"}, []string{"refuse prepare,the PrePrepare Message is not signed by the Leader Node"})
		return
	}
	p := Prepare{}
	s.SequenceId++
	s.MessagePool[pp.Digest] = pp.RequestMessage
	p.Digest = pp.Digest
	p.NodeId = s.Id
	sig, err = s.priKey.Sign(hash[:])
	if err != nil {
		utils.LogError([]string{"handlePrePrepare"}, []string{"refuse prepare,the PrePrepare Message is not signed by the Leader Node"})
	}

}
func (s *NodeServer) handlePrepare(payload []byte) {

}

func (s *NodeServer) Broadcast(rpc RPC) {
	data, err := json.Marshal(&rpc)
	if err != nil {
		utils.LogMsg([]string{"broadcast"}, []string{"rpc marshal error"})
	}
	for _, node := range NodeTables {
		if node.Id == s.Node.Id {
			continue
		}
		tcpDial(data, node.Addr)
	}
	utils.LogMsg([]string{"broadcast"}, []string{"broadcast completed"})
}

func randomId() int {
	num := 0
	// 获取一个 10 位数的 id
	for num < 1000000000 {
		num = rand.Intn(10000000000)
	}
	return num
}
