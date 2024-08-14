package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	// whether the height is correct
	if b.Height != int32(v.bc.Height()+1) {
		return fmt.Errorf("block with height %d is not the next block to be added, the actual height should be %d", b.Height, v.bc.Height()+1)
	}
	// whether the block is verified
	if err := b.Verify(); err != nil {
		return err
	}
	prevBlock, err := v.bc.GetBlock(int(b.Height) - 1)
	if err != nil {
		return nil
	}
	hash := BlockHasher{}.Hash(prevBlock.Header)
	// whether the prevHash is correct
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block is %v", b.PrevBlockHash)
	}
	return nil
}
