package network

import (
	"agricultural_meta/core"
	"agricultural_meta/types"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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
