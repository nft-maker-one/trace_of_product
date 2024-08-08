package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	priKey := GeneratePrivateKey()
	pubKey := priKey.PublicKey()
	address := pubKey.Address()
	fmt.Println(address.String())
	msg := []byte("hello world")
	msgError := []byte("hi world")
	sig, err := priKey.Sign(msg)
	assert.Nil(t, err)
	fmt.Println(sig)
	assert.True(t, sig.Verify(pubKey, msg))
	assert.False(t, sig.Verify(pubKey, msgError))
}
