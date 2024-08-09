package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Header struct {
	Version       int32
	PrevBlockHash types.Hash
	DataHash      types.Hash
	Timestamp     int64
	Height        int32
	Nonce         int64
}

type Block struct {
	*Header
	Eggplants []*Eggplant
	Validator crypto.PublicKey
	Signature *crypto.Signature
	hash      types.Hash
}

func NewBlock(h *Header, eggs []*Eggplant) (*Block, error) {
	return &Block{Header: h, Eggplants: eggs}, nil
}

func (b *Block) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(b.HeaderData())
	if err != nil {
		return err
	}
	b.Signature = sig
	b.Validator = priKey.PublicKey()
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	dataHash, err := CalculateDataHash(b.Eggplants)
	if err != nil {
		return err
	}
	if b.Header.DataHash != dataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}
	return nil
}

func (b *Block) AddTransaction(egg *Eggplant) {
	b.Eggplants = append(b.Eggplants, egg)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)
	logrus.WithField("headerData", buf.Bytes()).Debugln()
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

func CalculateDataHash(eggs []*Eggplant) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}
	for _, egg := range eggs {
		if err := NewGobEggplantEncoder(buf).Encode(egg); err != nil {
			return hash, err
		}

	}

	return types.Hash(sha256.Sum256(buf.Bytes())), nil
}

func NewBlockFromPrevHeader(prevHeader *Header, eggs []*Eggplant) (*Block, error) {
	dataHash, err := CalculateDataHash(eggs)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version:       1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
		Height:        prevHeader.Height + 1,
	}
	return NewBlock(header, eggs)

}
