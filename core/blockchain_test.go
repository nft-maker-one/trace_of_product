package core

import (
	"agricultural_meta/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(nil, randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc

}

func TestAddBlock(t *testing.T) {
	bc := newBlockWithGenesis(t)
	for i := 1; i <= 1000; i++ {
		block := randomBlock(t, i, getPrevBlockHash(t, bc, i))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, bc.Height(), 1000)
	assert.Equal(t, len(bc.Headers), 1001)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 88, getPrevBlockHash(t, bc, 88))))
	assert.Nil(t, bc.AddBlock(randomBlock(t, 1001, getPrevBlockHash(t, bc, 1001))))

}

func TestBlockchain(t *testing.T) {
	bc := newBlockWithGenesis(t)

	assert.NotNil(t, bc.validator)
	fmt.Println(bc.Height())
}

func TestHasBlock(t *testing.T) {
	bc := newBlockWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockWithGenesis(t)
	for i := 1; i <= 1000; i++ {
		b := randomBlock(t, i, getPrevBlockHash(t, bc, i))
		bc.AddBlock(b)
		header, _ := bc.GetHeader(int(b.Height))
		assert.Equal(t, b.Header, header)

	}

}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height int) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
