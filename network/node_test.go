package network

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {
	node1 := NewNodeServer("127.0.0.1:8081")
	node2 := NewNodeServer("127.0.0.1:8082")
	node3 := NewNodeServer("127.0.0.1:8083")
	node4 := NewNodeServer("127.0.0.1:8084")
	NodeTables = append(NodeTables, node1.Node)
	NodeTables = append(NodeTables, node2.Node)
	NodeTables = append(NodeTables, node3.Node)
	NodeTables = append(NodeTables, node4.Node)
	go node1.NodeUp()
	go node2.NodeUp()
	go node3.NodeUp()
	go node4.NodeUp()
	message := Message{}
	message.Header = MessageTypeTest
	message.Data = []byte("hello world!!!")
	buf := &bytes.Buffer{}
	assert.Nil(t, gob.NewEncoder(buf).Encode(&message))
	rpc := RPC{}
	rpc.ContentType = CommandType("test")
	rpc.Payload = buf.Bytes()
	node1.Broadcast(rpc)
	time.Sleep(5 * time.Second)
}
