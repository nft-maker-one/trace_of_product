package network

import (
	"agricultural_meta/core"
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTransForm MessageType = 0x1
	MessageTypeProcess   MessageType = 0x2
	MessageTypeStorage   MessageType = 0x3
	MessageTypeSell      MessageType = 0x4
	MessageTypeTx        MessageType = 0x5
	MessageTypeBlock     MessageType = 0x6
)

type RPC struct {
	From    NetAddr
	Payload []byte
}

type RPCHandler interface {
	HandleRPC(rpc RPC) error
}

type Message struct {
	Header MessageType
	Data   []byte
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

func DefaultRPCDecodeFunc(rpc RPC) (*DecodeMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message")
	}
	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")
	switch msg.Header {
	case MessageTypeTx:
		egg := new(core.Eggplant)
		if err := egg.Decode(core.NewGobEggplantDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: egg,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecode(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: block,
		}, nil
	default:
		return nil, nil
	}
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCProcesser {
	return &DefaultRPCProcesser{
		p: p,
	}
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{Header: t, Data: data}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

func (h *DefaultRPCProcesser) HandleRPC(rpc RPC) (*DecodeMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s : %s", rpc.From, err.Error())
	}
	switch msg.Header {
	case MessageTypeTx:
		tx := core.NewEggplant([]byte("initial"))
		if err := tx.Decode(core.NewGobEggplantDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message type %x", msg.Header)
	}
}
