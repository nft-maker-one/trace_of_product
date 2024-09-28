package crypto

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	for i := 0; i < 1000; i++ {
		priKey := GeneratePrivateKey()
		pubKey := priKey.PublicKey().ToSlice()
		msg := []byte("hello world" + strconv.Itoa(i))
		msgError := []byte("hi world" + strconv.Itoa(i))
		sig, err := priKey.Sign(msg)
		assert.Nil(t, err)
		assert.True(t, sig.Verify(pubKey, msg))
		assert.False(t, sig.Verify(pubKey, msgError))
	}

}

func TestGenerateKeyPair(t *testing.T) {
	test_data := ""
	for i := 0; i < 10; i++ {
		num := rand.Intn(500)
		key := GenerateKeyPair(num)
		test_data += fmt.Sprintf("key id %d\n", num)
		assert.NotNil(t, key)
	}
	fmt.Println(test_data)

}

func TestReadKeyFromFile(t *testing.T) {
	ids := []int{450, 440, 176, 86, 307, 234, 195, 178, 147, 393}
	for _, id := range ids {
		priKey, err := ReadPriKey(id)
		assert.Nil(t, err)
		pubKey, err := ReadPubKey(id)

		assert.Nil(t, err)
		msg := []byte("hello jerry")
		sig, err := priKey.Sign(msg)
		assert.Nil(t, err)
		assert.True(t, sig.Verify((*pubKey).ToSlice(), msg))
	}
}

func TestSignatureToByte(t *testing.T) {
	for i := 0; i < 1000; i++ {
		priKey := GeneratePrivateKey()
		pubKey := priKey.PublicKey()
		data := []byte("hello jerry")
		sig, err := priKey.Sign(data)
		assert.Nil(t, err)
		sigByte := sig.ToByte()
		sig1, err := ByteToSignature(sigByte)
		assert.Nil(t, err)
		assert.True(t, sig1.Verify(pubKey.ToSlice(), data))
		assert.Equal(t, sig, sig1)
	}

}
