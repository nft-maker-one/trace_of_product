package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"encoding/gob"
	"fmt"
)

type MetaData struct {
	EggplantId      int
	ProductHeight   int
	ProductHash     types.Hash
	TransportHeight int
	TransportHash   types.Hash
	ProcessHeight   int
	ProcessHash     types.Hash
	StorageHeight   int
	StorageHash     types.Hash
	SellHeight      int
	SellHash        types.Hash
}

type Eggplant struct {
	MetaData
	PublickKey crypto.PublicKey
	Signature  *crypto.Signature
	hash       types.Hash
	firstSeen  int64
}

func NewEggplant(data []byte) *Eggplant {
	return &Eggplant{PublickKey: crypto.PublicKey{}, Signature: &crypto.Signature{}}
}

func (eg *Eggplant) SetFirsstSeen(t int64) {
	eg.firstSeen = t
}

func (eg *Eggplant) FirstSeen() int64 {
	return eg.firstSeen
}

func (eg *Eggplant) SetHash(hash types.Hash) {
	eg.hash = hash
}

func (eg *Eggplant) EncodeMetaData() ([]byte, error) {
	buf := &bytes.Buffer{}
	if eg.EggplantId < 0 {
		return nil, fmt.Errorf("no valid MetaData")
	}
	if err := gob.NewEncoder(buf).Encode(eg.MetaData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (eg *Eggplant) Sign(priKey crypto.PrivateKey) error {
	data, err := eg.EncodeMetaData()
	if err != nil {
		return err
	}
	sig, err := priKey.Sign(data)
	if err != nil {
		return err
	}
	eg.PublickKey = priKey.PublicKey()
	eg.Signature = sig
	return nil
}

func (eg *Eggplant) Verify() error {
	data, err := eg.EncodeMetaData()
	if err != nil {
		return err
	}
	if eg.Signature == nil {
		return fmt.Errorf("Invalid Signature")
	}
	if !eg.Signature.Verify(eg.PublickKey, data) {
		return fmt.Errorf("Invalid Signature")
	}
	return nil
}

func (eg *Eggplant) Hash(hasher Hasher[*Eggplant]) types.Hash {
	if eg.hash.IsZero() {
		eg.hash = hasher.Hash(eg)
	}
	return hasher.Hash(eg)
}

func (eg *Eggplant) Decode(dec Decoder[*Eggplant]) error {
	return dec.Decode(eg)
}

func (eg *Eggplant) Encode(enc Encoder[*Eggplant]) error {
	return enc.Encode(eg)
}
