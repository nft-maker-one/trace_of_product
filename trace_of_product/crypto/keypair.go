package crypto

import (
	"agricultural_meta/types"
	"agricultural_meta/utils"
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
	key.Curve = elliptic.P256()
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

func (sig *Signature) Verify(pubKey []byte, data []byte) bool {
	pubkey := ByteToPubKey(pubKey)
	if pubkey.Key.X == nil {
		utils.LogMsg([]string{"Verify"}, []string{"public key with wrong format"})
		return false
	}
	res := ecdsa.Verify(pubkey.Key, data, sig.R, sig.S)
	return res
}

// 签名编码模块
func (sig *Signature) ToByte() []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	compressed := make([]byte, 2+len(rBytes)+len(sBytes))
	compressed[0] = byte(len(rBytes))
	compressed[1] = byte(len(sBytes))
	copy(compressed[2:2+len(rBytes)], rBytes)
	copy(compressed[2+len(rBytes):], sBytes)
	return compressed
}

// 签名解码模块
func ByteToSignature(compressed []byte) (*Signature, error) {
	rLen, sLen := int(compressed[0]), int(compressed[1])
	if len(compressed) != 2+sLen+rLen {
		return nil, fmt.Errorf("sig encoded format error")
	}
	r, s := big.NewInt(0), big.NewInt(0)
	r.SetBytes(compressed[2 : 2+rLen])
	s.SetBytes(compressed[2+rLen:])
	return &Signature{R: r, S: s}, nil
}
