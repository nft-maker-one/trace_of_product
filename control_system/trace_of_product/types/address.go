package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSilce() []byte {
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}
	return b
}

func NewAddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20\n", len(b))
		panic(msg)
	}
	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}
	return Address(value)
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSilce())
}
