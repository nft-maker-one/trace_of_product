package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"fmt"
	"testing"
	"time"

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
	test_data := ""
	res := ""
	for i := 0; i < 5; i++ {
		block := randomBlock(t, 1, types.Hash{})
		hash := block.Hash(BlockHasher{})
		test_data += fmt.Sprintf("block %d\n%+v\n", i+1, block)
		res += fmt.Sprintf("hash %d : %v\n", i+1, hash)
		assert.False(t, hash.IsZero())

	}
	fmt.Println(test_data)
	fmt.Println(res)

}

func TestBlockVerify(t *testing.T) {
	test_data := ""
	for i := 0; i < 5; i++ {
		b := randomBlock(t, 0, types.Hash{})
		priKey := crypto.GeneratePrivateKey()
		assert.Nil(t, b.Sign(priKey))
		assert.NotNil(t, b.Signature)
		test_data += fmt.Sprintf("block %d\n%+v\n", i+1, b)
		assert.Nil(t, b.Verify())
	}
	fmt.Println(test_data)
}

func TestEncodeBlock(t *testing.T) {
	test_data := ""
	for i := 0; i < 5; i++ {
		b := randomBlock(t, 1, types.Hash{})
		priKey := crypto.GeneratePrivateKey()
		b.Sign(priKey)
		test_data += fmt.Sprintf("block %d\n%+v\n", i+1, b)
		buf := bytes.Buffer{}
		assert.Nil(t, b.Encode(NewGobBlockEncoder(&buf)))
		bDecode := new(Block)
		assert.Nil(t, bDecode.Decode(NewGobBlockDecode(&buf)))
		assert.Equal(t, b, bDecode)
	}
	fmt.Println(test_data)
}
