package network

import (
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"crypto/sha256"
	"encoding/json"
)

type MessageType byte

const (
	MessageTypeProduce   MessageType = 0x1
	MessageTypeTransport MessageType = 0x2
	MessageTypeProcess   MessageType = 0x3
	MessageTypeStorage   MessageType = 0x4
	MessageTypeSell      MessageType = 0x5
	MessageTypeBlock     MessageType = 0x6
	MessageTypeTest      MessageType = 0x7
)

type RPC struct {
	ContentType CommandType `json:"content_type"`
	Payload     []byte      `json:"payload"`
}

type RPCHandler interface {
	HandleRPC(rpc RPC) error
}

type Message struct {
	Header MessageType `json:"header"`
	Data   []byte      `json:"data"`
}

type RPCProcessor interface {
	ProcessMessage(*DecodeMessage) error
}

type DefaultRPCProcesser struct {
	p RPCProcessor
}

type DecodeMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCProcesser {
	return &DefaultRPCProcesser{
		p: p,
	}
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{Header: t, Data: data}
}

func (m Message) Hash() types.Hash {
	data, err := json.Marshal(&m)
	if err != nil {
		utils.LogMsg([]string{"Message.Hash"}, []string{"marshal failed"})
		return types.Hash{}
	}
	hash := sha256.Sum256(data)
	return types.Hash(hash)
}
