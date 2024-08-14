package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) Connect(tr Transport) error {
	// register the address of other nodes
	t.lock.Lock()
	defer t.lock.Unlock()
	ltr, ok := tr.(*LocalTransport)
	if !ok {
		return fmt.Errorf("input object is not in a structure of LocalTransport")
	}
	t.peers[tr.Addr()] = ltr
	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	// payload 是接受的序列化后的 Message 信息
	// send message to all the nodes registered in its peers list.
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	// payload 是序列化后的 Message 信息
	// send the message to a specific node
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.addr == to {
		return nil
	}
	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s could not send message to %s", t.Addr(), to)
	}
	// payload 是对 msg 的序列化
	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}
	return nil
}
