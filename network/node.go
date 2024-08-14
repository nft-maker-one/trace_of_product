package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/types"
)

type Node struct {
	Id   int    `gorm:"id"`
	Addr string `gorm:"addr"`
}

type NodeServer struct {
	Chain               *core.Blockchain
	messagePool         map[types.Hash]Request
	prePareConfirmCount map[types.Hash]map[string]int
	commitConfirmCount  map[types.Hash]map[string]int
	isCommitBroadcast   map[types.Hash]bool
	isReply             map[types.Hash]bool
	priKey              *crypto.PrivateKey
	PubKey              *crypto.PublicKey
}

func NewNode() (*Node, error) {
	return nil, nil
}
