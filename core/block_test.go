package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, height int, prevBlockHash types.Hash) *Block {
	return nil
}

func TestHashBlock(t *testing.T) {

}

func TestBlockVerify(t *testing.T) {
	b := randomBlock(t, 0, types.Hash{})
	priKey := crypto.GeneratePrivateKey()
	assert.Nil(t, b.Sign(priKey))
	assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())
	b.Height = 103
	assert.NotNil(t, b.Verify())
}

func TestEncodeBlock(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	buf := bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(&buf)))
	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecode(&buf)))
	assert.Equal(t, b, bDecode)
}
