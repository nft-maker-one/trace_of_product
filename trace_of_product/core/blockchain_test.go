package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc

}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height int) types.Hash {
	prevBlock, err := bc.GetBlock(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevBlock.Header)
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	for i := 1; i <= 1000; i++ {
		block := randomBlock(t, i, getPrevBlockHash(t, bc, i))
		priKey := crypto.GeneratePrivateKey()
		assert.Nil(t, block.Sign(priKey))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, bc.Height(), 1000)
	assert.Equal(t, len(bc.Chains), 1001)
	block1 := randomBlock(t, 88, getPrevBlockHash(t, bc, 88))
	prikey1 := crypto.GeneratePrivateKey()
	block1.Sign(prikey1)
	prikey2 := crypto.GeneratePrivateKey()
	block2 := randomBlock(t, 1001, getPrevBlockHash(t, bc, 1001))
	block2.Sign(prikey2)
	assert.NotNil(t, bc.AddBlock(block1))
	assert.Nil(t, bc.AddBlock(block2))

}

func TestBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	fmt.Println(bc.Height())
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t) //block Height 0 Chains Length 1
	fmt.Println(bc.Height())
	fmt.Println(bc.GetBlock(0))
	for i := 1; i <= 1000; i++ {
		b := randomBlock(t, i, getPrevBlockHash(t, bc, i)) //
		priKey := crypto.GeneratePrivateKey()
		b.Sign(priKey)
		err := bc.AddBlock(b)
		assert.Nil(t, err)
		block, err := bc.GetBlock(int(b.Height))
		fmt.Println(bc.Height())
		assert.Nil(t, err)
		assert.Equal(t, b.Header, block.Header)

	}

}
