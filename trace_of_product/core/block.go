package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

const AddPoint = 1

type Header struct {
	Version       int32
	PrevBlockHash types.Hash
	DataHash      types.Hash
	Timestamp     int64
	Height        int32
	Nonce         int64
	Scores        map[int]int
	Leader        int
}

// Datahash 是记录对区块中茄子信息的哈希
// hash 是整个区块信息的哈希
type Block struct {
	*Header
	Eggplants []*Eggplant
	Validator []byte
	Signature []byte
	BlockHash types.Hash
}

func (h *Header) UpdateScore(bc *Blockchain, eggs []*Eggplant) {
	prevHeader := bc.Chains[bc.Height()]
	prevScore := prevHeader.Scores
	newScore := make(map[int]int)
	for _, egg := range eggs {
		if _, ok := newScore[egg.NodeId]; !ok {
			newScore[egg.NodeId] = prevScore[egg.NodeId] / 2
		}
		newScore[egg.EggplantId]++
	}
	prevScore[prevHeader.Leader] = 0
	h.Scores = prevScore
}

func (h *Header) SelectLeader() {
	max := 0
	leader := 0
	for k, v := range h.Scores {
		if v > max {
			leader = k
			max = v
		}
	}
	h.Leader = leader
}

func NewBlock(h *Header, eggs []*Eggplant) (*Block, error) {
	return &Block{Header: h, Eggplants: eggs}, nil
}

func (b *Block) Sign(priKey crypto.PrivateKey) error {
	// use ecdsa to sign the headerData
	sig, err := priKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Signature = sig.ToByte()
	b.Validator = priKey.PublicKey().ToSlice()
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no validator")
	}
	sig, err := crypto.ByteToSignature(b.Signature)
	if err != nil {
		return err
	}

	if !sig.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid validator")
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

func (b *Block) AddEggplant(egg *Eggplant) {
	b.Eggplants = append(b.Eggplants, egg)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.BlockHash.IsZero() {
		b.BlockHash = hasher.Hash(b.Header)
	}
	return b.BlockHash
}

func (b *Block) GetHash() types.Hash {
	if b.BlockHash.IsZero() {
		b.BlockHash = BlockHasher{}.Hash(b.Header)
	}
	return b.BlockHash
}

// use gobEncoder to encode Header Data
func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)
	return buf.Bytes()
}

func (h *Header) Bytes() []byte {
	// buf := &bytes.Buffer{}
	// enc := gob.NewEncoder(buf)
	// enc.Encode(h)
	// return buf.Bytes()
	res := make([]byte, 0)
	res = append(res, byte(h.Version))
	res = append(res, h.PrevBlockHash[:]...)
	res = append(res, h.DataHash[:]...)
	res = append(res, byte(h.Timestamp))
	res = append(res, byte(h.Height))
	res = append(res, byte(h.Nonce))
	for k, v := range h.Scores {
		res = append(res, byte(k+v))
	}
	res = append(res, byte(h.Leader))
	return res
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
