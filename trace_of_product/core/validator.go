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

	prevBlock := v.bc.Chains[v.bc.Height()]
	// whether the prevHash is correct
	if b.PrevBlockHash != prevBlock.BlockHash {
		return fmt.Errorf("the hash of the previous block is %v but acutual got %v", prevBlock.BlockHash, b.PrevBlockHash)
	}
	return nil
}
