package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	fmt.Println(tra)
	fmt.Println(trb.peers[tra.addr])
	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	m := Message{}
	rpc := <-trb.Consume()
	assert.Nil(t, gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&m))
	assert.Equal(t, m.Data, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("foo")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	msgb := Message{}
	assert.Nil(t, gob.NewDecoder(bytes.NewReader(rpcb.Payload)).Decode(&msgb))
	assert.Equal(t, msg, msgb.Data)
}
