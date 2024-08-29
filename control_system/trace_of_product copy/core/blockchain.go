package core

import (
	"agricultural_meta/types"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	lock      sync.RWMutex
	Chains    []*Block
	validator Validator
}

func CreateChain(leadrId int) *Blockchain {
	chain := &Blockchain{}
	chain.lock = sync.RWMutex{}
	chain.validator = NewBlockValidator(chain)
	chain.Chains = make([]*Block, 0)
	h := &Header{
		Version:       0,
		PrevBlockHash: types.Hash{},
		DataHash:      types.Hash{},
		Timestamp:     0,
		Height:        0,
		Nonce:         0,
		Leader:        leadrId,
		Scores:        make(map[int]int),
	}
	b, err := NewBlock(h, nil)
	if err != nil {
		panic(err)
	}
	err = chain.addBlockWithoutValidation(b)
	if err != nil {
		panic(err)
	}
	return chain

}

func (bc *Blockchain) Height() int {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return len(bc.Chains) - 1
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetBlock(height int) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("there is not block with height %d", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.Chains[height], nil
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) HasBlock(height int) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) GetLeader() int {
	lastBlock, err := bc.GetBlock(bc.Height())
	if err != nil {
		panic(fmt.Errorf("failed to get leader err=%s", err.Error()))
	}
	return lastBlock.Leader
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if bc.Chains == nil {
		bc.Chains = make([]*Block, 0)
	}
	logrus.WithFields(
		logrus.Fields{
			"msg":  "create block",
			"hash": b.Hash(BlockHasher{}),
		},
	).Println()
	bc.Chains = append(bc.Chains, b)

	return nil
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		lock:   sync.RWMutex{},
		Chains: make([]*Block, 0),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) GetPrevHash() types.Hash {
	return bc.Chains[bc.Height()].GetHash()
}
