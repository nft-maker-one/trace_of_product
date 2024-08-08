package network

import (
	"agricultural_meta/core"
	"agricultural_meta/types"
	"sort"
)

type TxPool struct {
	eggplants map[types.Hash]*core.Eggplant
}

type TxMapSorter struct {
	eggplants []*core.Eggplant
}

func NewMapSorter(txMap map[types.Hash]*core.Eggplant) *TxMapSorter {
	txx := make([]*core.Eggplant, len(txMap))
	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}
	return &TxMapSorter{eggplants: txx}
}

func (s *TxMapSorter) Len() int {
	return len(s.eggplants)
}

func (s *TxMapSorter) Swap(i, j int) {
	s.eggplants[i], s.eggplants[j] = s.eggplants[j], s.eggplants[i]
}

func (s *TxMapSorter) Less(i, j int) bool {
	return s.eggplants[i].FirstSeen() < s.eggplants[j].FirstSeen()
}

func NewTxPool() *TxPool {
	return &TxPool{
		eggplants: make(map[types.Hash]*core.Eggplant),
	}
}

func (p *TxPool) Transactions() []*core.Eggplant {
	s := NewMapSorter(p.eggplants)
	sort.Sort(s)
	return s.eggplants
}

func (p *TxPool) Len() int {
	return len(p.eggplants)
}

func (p *TxPool) Has(hash types.Hash) bool {
	if _, ok := p.eggplants[hash]; ok {
		return true
	}
	return false
}

func (p *TxPool) Add(tx *core.Eggplant) error {
	hash := tx.Hash(core.EggplantHasher{})
	if p.Has(hash) {
		return nil
	}
	p.eggplants[hash] = tx
	return nil
}

func (p *TxPool) Flush() {
	p.eggplants = make(map[types.Hash]*core.Eggplant)
}

// func (p *TxPool) Flush() {
// 	p.transactions = make(map[types.Hash]*core.Transaction)
// }
