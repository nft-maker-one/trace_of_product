package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("the Provided bytes with length %d should be 32", len(b))
		panic(msg)
	}
	var value [32]uint8
	for i := 0; i < 32; i++ {
		value[i] = b[i]
	}
	return Hash(value)
}

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() Hash {
	return Hash(RandomBytes(32))
}

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func (h Hash) ToSlice() []byte {
	res := make([]byte, 32)
	for i := 0; i < 32; i++ {
		res[i] = h[i]
	}
	return res
}
