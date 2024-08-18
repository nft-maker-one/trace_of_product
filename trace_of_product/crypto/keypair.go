package crypto

import (
	"agricultural_meta/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{r, s}, nil
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func ByteToPubKey(data []byte) PublicKey {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	key := ecdsa.PublicKey{}
	key.X = x
	key.Y = y
	return PublicKey{Key: &key}
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.NewAddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	res := ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
	return res
}

// 签名编码模块
func (sig *Signature) ToByte() []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	bytesLength := (elliptic.P256().Params().BitSize + 7) / 8
	compressed := make([]byte, 2*bytesLength)
	copy(compressed[:bytesLength], rBytes)
	copy(compressed[bytesLength:], sBytes)
	return compressed
}

// 签名解码模块
func ByteToSignature(compressed []byte) (*Signature, error) {
	if len(compressed)%2 != 0 {
		return nil, fmt.Errorf("bytes array has a invalid odd length")
	}
	byteLength := len(compressed) / 2
	if byteLength != (elliptic.P256().Params().BitSize+7)/8 {
		return nil, fmt.Errorf("bytes array is not a signature by secp256r1")
	}
	r := new(big.Int).SetBytes(compressed[:byteLength])
	s := new(big.Int).SetBytes(compressed[byteLength:])
	return &Signature{R: r, S: s}, nil
}
