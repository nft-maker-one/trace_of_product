package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type NodeType int

const X float64 = 0.2
const Y float64 = 0.3

const (
	Supervisor NodeType = iota
	Normal
	Down
)

type ServerOpts struct {
	ID            string
	Cred          map[*Server]float64
	Type          NodeType
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
	Blockchain    *core.Blockchain
}

type Server struct {
	ServerOpts
	chain       *core.Blockchain
	blockTime   time.Duration
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
	memPool     *TxPool
}

func NewServer(opts ServerOpts) (*Server, error) {

	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, geensisBlock())

	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
		memPool:     NewTxPool(),
		chain:       chain,
		isValidator: opts.PrivateKey != nil,
	}
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) Start() {
	s.initTransports()
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Println(string(rpc.From) + " =======> " + s.ID)
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}

		case <-s.quitCh:
			break free
		}
	}
	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc

			}

		}(tr)
	}
}

func (s *Server) ProcessNewEggplant(egg *core.Eggplant) error {

	hash := egg.Hash(core.EggplantHasher{})
	if s.memPool.Has(hash) {
		return nil
	}
	egg.SetHash(hash)
	if err := egg.Verify(); err != nil {
		return err
	}
	egg.SetFirsstSeen(time.Now().UnixNano())

	// s.Logger.Log("msg", "adding new tx to mempool", "hash", hash, "mempoolLength", s.memPool.Len())
	go s.broadcastTx(egg)
	return s.memPool.Add(egg)
}

func (s *Server) ProcessMessage(msg *DecodeMessage) error {

	switch msg.Data.(type) {
	case *core.Eggplant:
		return s.ProcessNewEggplant((msg.Data).(*core.Eggplant))
	case *core.Block:
		return s.ProcessBlock((msg.Data).(*core.Block))
	}
	return nil
}

func (s *Server) ProcessBlock(b *core.Block) error {
	err := s.chain.AddBlock(b)
	// if err != nil {
	// 	validator := s.FindNodeByValidator(b.Validator)
	// 	s.updateCred(validator, s.IsSupervisor(b), X, Y)

	// }
	return err
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			logrus.Error(err.Error())
			return err
		}
	}
	return nil
}

func (s *Server) broadcastTx(tx *core.Eggplant) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobEggplantEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

func geensisBlock() *core.Block {
	header := core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}
	b, _ := core.NewBlock(&header, nil)
	return b
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}
	txx := s.memPool.Transactions()
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}
	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}
	if err := s.chain.AddBlock(block); err != nil {
		return err
	}
	go s.broadcastBlock(block)
	s.memPool.Flush()
	return nil
}

func (s *Server) updateCred(node *Server, isSignature bool, x, y float64) {
	if isSignature {
		if node.Type == Supervisor {
			s.Cred[node] = (1 + x) * s.Cred[node]
		} else if node.Type == Normal {
			s.Cred[node] = (1 + y) * s.Cred[node]
		}
	} else {
		if node.Type == Supervisor {
			s.Cred[node] = maxI(0, (1-x)*s.Cred[node])
		} else if node.Type == Normal {
			s.Cred[node] = maxI(0, (1-y)*s.Cred[node])
		}
	}

}

func maxI(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (s *Server) FindNodeByValidator(validator crypto.PublicKey) *Server {
	return nil
}

func (s *Server) IsSupervisor(b *core.Block) bool {
	return true
}
