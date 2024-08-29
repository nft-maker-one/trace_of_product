package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, height int, prevBlockHash types.Hash) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Height:        int32(height),
		Nonce:         0,
	}
	eggs := []*Eggplant{randomEggplantWithSignature()}
	newBlock, err := NewBlock(h, eggs)
	assert.Nil(t, err)
	hash, err := CalculateDataHash(newBlock.Eggplants)
	assert.Nil(t, err)
	newBlock.DataHash = hash
	// block := NewBlock()
	return newBlock
}

func TestHashBlock(t *testing.T) {
	block := randomBlock(t, 1, types.Hash{})
	logrus.WithField("block", block).Debugln()
	hash := block.Hash(BlockHasher{})
	assert.False(t, hash.IsZero())
	logrus.WithField("hash", hash).Infoln()
}

func TestBlockVerify(t *testing.T) {
	b := randomBlock(t, 0, types.Hash{})
	priKey := crypto.GeneratePrivateKey()
	assert.Nil(t, b.Sign(priKey))
	assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())
}

func TestEncodeBlock(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	priKey := crypto.GeneratePrivateKey()
	b.Sign(priKey)
	buf := bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(&buf)))
	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecode(&buf)))
	assert.Equal(t, b, bDecode)
}
