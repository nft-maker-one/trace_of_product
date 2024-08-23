package network

import (
	"agricultural_meta/core"
	"agricultural_meta/types"
	"encoding/binary"
	"encoding/gob"
	"io"
)

type ReplyStatus int
type CommandType string

type Request struct {
	Message
	TImestamp int64
}

type PrePrepare struct {
	RequestMessage core.Block `json:"request_message"`
	Digest         types.Hash `json:"digest"`
	SequencId      int        `json:"sequence_id"`
	Sign           []byte     `json:"sign"`
}

type Prepare struct {
	Digest    types.Hash `json:"digest"`
	SequencId int        `json:"sequenc_id"`
	NodeId    int        `json:"node_id"`
	Sign      []byte     `json:"sign"`
}

type Commit struct {
	Digest     types.Hash `json:"digest"`
	SequenceId int        `json:"sequence_id"`
	NodeId     int        `json:"node_id"`
	Sign       []byte     `json:"sign"`
	Nonce      uint64
}

type Reply struct {
	NodeId  int    `json:"node_id"`
	Content string `json:"content"`
}

type EggMes struct {
	NodeId int
	Egg    []byte
}

const (
	Consensus ReplyStatus = 1
	Failed    ReplyStatus = 2
)

const (
	cRequest    CommandType = "request"
	cPrePrepare CommandType = "preprepare"
	cPrepare    CommandType = "prepare"
	cCommit     CommandType = "commit"
	cReply      CommandType = "reply"
	cEgg        CommandType = "egg"
	cTest       CommandType = "test"
)

func CommitData(hash types.Hash, nonce uint64) []byte {
	// 10000000000 -> 34 位 2 进制
	res := make([]byte, 66)
	for i := 0; i < 32; i++ {
		res[i] = hash[i]
	}
	binary.BigEndian.PutUint64(res[32:], nonce)
	return res
}

func CommitDataSplit(data []byte) (types.Hash, uint64) {
	hash := types.Hash{}
	for i := 0; i < 32; i++ {
		hash[i] = data[i]
	}
	nonce := binary.BigEndian.Uint64(data[32:])
	return hash, nonce
}

func (e *EggMes) Encode(enc core.Encoder[*EggMes]) error {
	return enc.Encode(e)
}

func (e *EggMes) Decode(dec core.Decoder[*EggMes]) error {
	return dec.Decode(e)
}

type EggMesEncoder struct {
	w io.Writer
}

func NewGobEggMesEncoder(w io.Writer) *EggMesEncoder {
	return &EggMesEncoder{
		w: w,
	}
}

func (enc *EggMesEncoder) Encode(e *EggMes) error {
	return gob.NewEncoder(enc.w).Encode(e)
}

type EggMesDecoder struct {
	r io.Reader
}

func NewGobEggMesDecoder(r io.Reader) *EggMesDecoder {
	return &EggMesDecoder{
		r: r,
	}
}

func (dec *EggMesDecoder) Decode(e *EggMes) error {
	return gob.NewDecoder(dec.r).Decode(e)
}
