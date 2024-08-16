package network

import "agricultural_meta/types"

type ReplyStatus int
type CommandType string

type Request struct {
	Message
	TImestamp int64
}

type PrePrepare struct {
	RequestMessage Message    `json:"request_message"`
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
}

type Reply struct {
	Status  ReplyStatus `json:"status"`
	Content string      `json:"content"`
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
	cTest       CommandType = "test"
)
