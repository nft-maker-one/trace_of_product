package network

import (
	"agricultural_meta/core"
	"fmt"
	"strconv"
	"testing"

	"math/rand"

	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	txPool := NewTxPool()
	assert.Equal(t, txPool.Len(), 0)
}

func TestTxPoolAdd(t *testing.T) {
	txPool := NewTxPool()
	tx := core.NewEggplant([]byte("fooo"))
	assert.Nil(t, txPool.Add(tx))
	assert.Equal(t, txPool.Len(), 1)
}

func TestSortTransaction(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewEggplant([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirsstSeen(int64(rand.Intn(10000 * 1000)))
		assert.Nil(t, p.Add(tx))
	}
	assert.Equal(t, txLen, p.Len())
	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
	for i := 0; i < 10; i++ {
		fmt.Println(txx[i])
	}

}
