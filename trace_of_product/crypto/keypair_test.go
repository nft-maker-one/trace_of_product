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

func TestGenerateKeyPair(t *testing.T) {
	key := GenerateKeyPair(1)
	assert.NotNil(t, key)
	fmt.Println(key)
}

func TestReadKeyFromFile(t *testing.T) {
	priKey, err := ReadPriKey(1)
	assert.Nil(t, err)
	pubKey, err := ReadPubKey(1)
	assert.Nil(t, err)
	msg := []byte("hello jerry")
	sig, err := priKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(*pubKey, msg))
}

func TestSignatureToByte(t *testing.T) {
	priKey := GeneratePrivateKey()
	pubKey := priKey.PublicKey()
	data := []byte("hello jerry")
	sig, err := priKey.Sign(data)
	assert.Nil(t, err)
	sigByte := sig.ToByte()
	sig1, err := ByteToSignature(sigByte)
	assert.Nil(t, err)
	assert.True(t, sig1.Verify(pubKey, data))
	assert.Equal(t, sig, sig1)
}
