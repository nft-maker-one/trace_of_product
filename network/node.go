package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"math/rand"
)

type Node struct {
	Id   int    `gorm:"id"`
	Addr string `gorm:"addr"`
}

type NodeServer struct {
	Node
	Chain               *core.Blockchain
	MessagePool         map[types.Hash]Request
	PrePareConfirmCount map[types.Hash]map[string]int
	CommitConfirmCount  map[types.Hash]map[string]int
	IsCommitBroadcast   map[types.Hash]bool
	isReply             map[types.Hash]bool
	priKey              *crypto.PrivateKey
	PubKey              crypto.PublicKey
	ConsumeCh           []RPC
}

func NewNodeServer(addr string) *NodeServer {
	id := randomId()
	server := &NodeServer{}
	server.Id = id
	server.Addr = addr
	server.PubKey = crypto.GenerateKeyPair(id)
	server.MessagePool = make(map[types.Hash]Request)
	server.PrePareConfirmCount = make(map[types.Hash]map[string]int)
	server.CommitConfirmCount = make(map[types.Hash]map[string]int)
	server.isReply = make(map[types.Hash]bool)
	var err error
	server.priKey, err = crypto.ReadPriKey(id)
	if err != nil {
		panic(err)
	}
	utils.LogMsg([]string{"NewNodeServer"}, []string{"create server successfully"})
	return server

}

func (s *NodeServer) handleRequest(data []byte) {

}

func randomId() int {
	num := 0
	// 获取一个 10 位数的 id
	for num < 1000000000 {
		num = rand.Intn(10000000000)
	}
	return num
}
