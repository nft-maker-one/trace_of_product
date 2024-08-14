package network

type ReplyStatus int
type CommandType string

type Request struct {
	Message
	TImestamp int64
}

type PrePrepare struct {
	RequestMessage Request
	Digest         string
	SequencId      int
	Sign           []byte
}

type Prepare struct {
	Digest    string
	SequencId int
	NodeId    int
	Sign      []byte
}

type Commit struct {
	Digest     string
	SequenceId int
	NodeId     int
	Sign       []byte
}

type Reply struct {
	Status  ReplyStatus
	Content string
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
)
