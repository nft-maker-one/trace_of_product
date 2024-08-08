package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	logger    log.Logger
	store     Storage
	lock      sync.RWMutex
	Headers   ([]*Header)
	validator Validator
}

func (bc *Blockchain) Height() int {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return len(bc.Headers) - 1
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height int) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("there is not block with height %d", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.Headers[height], nil
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		Headers: []*Header{},
		store:   &MemoryStore{},
		logger:  l,
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
	if bc.Headers == nil {
		bc.Headers = make([]*Header, 0)
	}
	// fmt.Println()
	// fmt.Println("===================")
	// bc.logger.Log("msg", "showing block")
	// fmt.Printf("%+v\n", b)
	// hash := b.Hash(BlockHasher{})
	// fmt.Println(hash)
	// fmt.Println("===================")
	// fmt.Println()

	bc.Headers = append(bc.Headers, b.Header)
	bc.logger.Log(
		"msg", "add New Block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Eggplants),
	)
	return bc.store.Put(b)
}
