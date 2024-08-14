package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"encoding/gob"
	"fmt"
)

// MetaData represents for the sotrage message and identical message of Eggplants
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

// Eggplant is maingly consist of 5 part
type Eggplant struct {
	MetaData
	PublickKey crypto.PublicKey  //the Validator of this Eggplant
	Signature  *crypto.Signature //the Signature of the Validator
	hash       types.Hash        //the digest for eggplant's metadata
	firstSeen  int64             //the creation time of the eggplant
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
	// check whether the id of eggplant is valid
	if eg.EggplantId < 0 {
		return nil, fmt.Errorf("no valid MetaData")
	}
	// use the gobEncoder to encode MetaData of eggplant
	if err := gob.NewEncoder(buf).Encode(eg.MetaData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (eg *Eggplant) Sign(priKey crypto.PrivateKey) error {
	// encode the MetaData
	data, err := eg.EncodeMetaData()
	if err != nil {
		return err
	}
	// use ECDSA to sign the MetaData
	sig, err := priKey.Sign(data)
	if err != nil {
		return err
	}
	// set the Validator for the MetaData
	eg.PublickKey = priKey.PublicKey()
	// record the signature
	eg.Signature = sig
	return nil
}

func (eg *Eggplant) Verify() error {
	// encode the MetaData
	data, err := eg.EncodeMetaData()
	if err != nil {
		return err
	}

	if eg.Signature == nil {
		return fmt.Errorf("Eggplant without signature")
	}
	if !eg.Signature.Verify(eg.PublickKey, data) {
		return fmt.Errorf("Eggplant with wrong validator")
	}
	return nil
}

func (eg *Eggplant) Hash(hasher Hasher[*Eggplant]) types.Hash {
	//whether the hash existS
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
