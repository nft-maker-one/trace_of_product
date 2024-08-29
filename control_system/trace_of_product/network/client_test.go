package network

import (
	"agricultural_meta/core"
	"agricultural_meta/types"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPbft(t *testing.T) {
	h := &core.Header{
		Version:       0,
		PrevBlockHash: types.RandomHash(),
		DataHash:      types.RandomHash(),
		Timestamp:     time.Now().Unix(),
		Height:        0,
		Nonce:         rand.Int63n(math.MaxInt64),
	}
	b, err := core.NewBlock(h, nil)
	assert.Nil(t, err)
	fmt.Println(b.Hash(core.BlockHasher{}).String())
	fmt.Println(b.Hash(core.BlockHasher{}).String())
	pp := PrePrepare{
		RequestMessage: *b,
	}
	data, err := json.Marshal(&pp)
	assert.Nil(t, err)
	ppNew := PrePrepare{}
	assert.Nil(t, json.Unmarshal(data, &ppNew))
	fmt.Println(ppNew.RequestMessage.Hash(core.BlockHasher{}).String())
}

func TestParseRequest(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:8081")
		assert.Nil(t, err)
		for {
			conn, _ := listener.Accept()
			method, path := handleConnection(conn)
			assert.Equal(t, method, "GET")
			assert.Equal(t, path, "/hello")
		}
	}()
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/hello", nil)
	assert.Nil(t, err)
	client := http.Client{}
	client.Do(req)
	time.Sleep(1 * time.Second)
}

func TestPost(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:8081")
		assert.Nil(t, err)
		for {
			conn, _ := listener.Accept()
			method, path := handleConnection(conn)
			assert.Equal(t, method, "POST")
			assert.Equal(t, path, "/hello")
		}
	}()
	data := []byte(`{"name":"zhangshan"}`)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/hello", bytes.NewReader(data))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	client.Do(req)
	time.Sleep(1 * time.Second)
}
