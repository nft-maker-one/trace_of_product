package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store     Storage
	lock      sync.RWMutex
	Chains    []*Block
	validator Validator
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

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		Chains: []*Block{},
		store:  &MemoryStore{},
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) HasBlock(height int) bool {
	return height <= bc.Height()
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
			"hash": b.hash,
		},
	).Println()
	bc.Chains = append(bc.Chains, b)

	return nil
}
