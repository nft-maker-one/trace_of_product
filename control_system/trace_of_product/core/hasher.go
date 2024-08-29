package core

import (
	"agricultural_meta/types"
	"crypto/sha256"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type EggplantHasher struct {
}

func (EggplantHasher) Hash(e *Eggplant) types.Hash {
	data, _ := e.EncodeMetaData()

	h := sha256.Sum256(data)
	return types.Hash(h)
}
