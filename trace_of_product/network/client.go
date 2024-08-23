package network

import (
	"agricultural_meta/core"
	"agricultural_meta/database"
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"time"
)

type ClientNode struct {
	addr      string
	db        *database.NodeDb
	ConsumeCh []chan RPC
}

func NewClient(addr string) *ClientNode {
	node := ClientNode{addr: addr}
	cfg, err := utils.NewConfig("./config.yaml")
	if err != nil {
		panic(err)
	}
	node.db = database.InitNodeDb(cfg.Mysql.Dsn)

	return &node
}

func (c *ClientNode) ClientUp() {
	go func() {
		listen, err := net.Listen("tcp", c.addr)
		if err != nil {
			panic(err)
		}
		for {
			conn, err := listen.Accept()
			if err != nil {
				utils.LogMsg([]string{"NewClient"}, []string{"listening error err = " + err.Error()})
			}
			data, err := io.ReadAll(conn)
			if err != nil {
				utils.LogMsg([]string{"NewClient"}, []string{"read conn err = " + err.Error()})
			}
			rpc := new(RPC)
			err = json.Unmarshal(data, rpc)
			if err != nil {
				utils.LogMsg([]string{"NewClient"}, []string{"not the format of RPC err = " + err.Error()})
			}
			switch rpc.ContentType {
			case cReply:
				message := new(Message)
				err = json.Unmarshal(rpc.Payload, message)
				if err != nil {
					utils.LogMsg([]string{"ClientUp"}, []string{})
				}
			default:
				fmt.Println("content-type not correct")
			}
		}
	}()
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请依据下面提升输入指令")
		fmt.Println("1. 更新茄子信息")
		cmd, err := input.ReadString('\n')
		if err != nil {
			utils.LogMsg([]string{"ClientUp"}, []string{"read from cmd failed err=" + err.Error()})
			continue
		}
		if cmd == "1" {
			data := core.MetaData{
				EggplantId:      rand.Intn(1000),
				ProductHeight:   rand.Intn(math.MaxInt),
				ProductHash:     types.RandomHash(),
				TransportHeight: rand.Intn(math.MaxInt),
				TransportHash:   types.RandomHash(),
				ProcessHeight:   rand.Intn(math.MaxInt),
				ProcessHash:     types.RandomHash(),
				StorageHeight:   rand.Intn(math.MaxInt),
				StorageHash:     types.RandomHash(),
				SellHeight:      rand.Intn(math.MaxInt),
				SellHash:        types.RandomHash(),
			}
			egg := core.NewEggplant(data)
			buf := &bytes.Buffer{}
			err = egg.Encode(core.NewGobEggplantEncoder(buf))
			if err != nil {
				utils.LogMsg([]string{"ClientUp"}, []string{"encode block eggplant err = " + err.Error()})
				continue
			}

			rpc := RPC{}
			rpc.Payload = buf.Bytes()
			rpc.ContentType = cEgg
			rByte, err := json.Marshal(rpc)
			if err != nil {
				utils.LogMsg([]string{"ClientUp"}, []string{"encode rpc error err = " + err.Error()})
				continue
			}
			tcpDial(rByte, "127.0.0.1:8082")

		}
	}
}

func ClientBlock(prevBlockHash types.Hash, height int) *core.Block {
	metaData := core.MetaData{
		EggplantId:      rand.Intn(math.MaxInt),
		ProductHeight:   rand.Intn(math.MaxInt),
		ProductHash:     types.RandomHash(),
		TransportHeight: rand.Intn(math.MaxInt),
		TransportHash:   types.RandomHash(),
		ProcessHeight:   rand.Intn(math.MaxInt),
		ProcessHash:     types.RandomHash(),
		StorageHeight:   rand.Intn(math.MaxInt),
		StorageHash:     types.RandomHash(),
		SellHeight:      rand.Intn(math.MaxInt),
		SellHash:        types.RandomHash(),
	}
	egg := &core.Eggplant{}
	egg.MetaData = metaData
	// egg.FirstSeen = time.Now().Unix()
	// priKey := crypto.GeneratePrivateKey()
	// if err := egg.Sign(priKey); err != nil {
	// 	logrus.Error(err)
	// 	return nil
	// }
	h := &core.Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Height:        int32(height),
		Nonce:         0,
	}
	block := new(core.Block)
	block.Header = h
	block.Eggplants = append(block.Eggplants, egg)
	hash, _ := core.CalculateDataHash(block.Eggplants)
	block.DataHash = hash
	return block

}
