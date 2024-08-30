package network

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoblyLinkedList(t *testing.T) {
	head := NewHead()
	hash1 := types.RandomHash()
	node1 := NewListNode(hash1)
	head.TailAdd(node1)
	assert.Equal(t, head.len, 1)
	hash2 := types.RandomHash()
	node2 := NewListNode(hash2)
	head.TailAdd(node2)
	hash3 := types.RandomHash()
	node3 := NewListNode(hash3)
	head.TailAdd(node3)
	assert.Equal(t, 3, head.len)
	popNode, err := head.TopPop()
	assert.Nil(t, err)
	assert.Equal(t, node1.value, popNode)
	assert.Equal(t, 2, head.len)
	assert.Equal(t, node3.value, head.prev.value)
	head.DeleteNode(node3)
	assert.Equal(t, head.prev.value, node2.value)
}

func TestAddEgg(t *testing.T) {
	pool := NewMemoryPool(500)
	for i := 1; i <= 100; i++ {
		data := core.MetaData{EggplantId: i}
		egg := core.NewEggplant(data)
		priKey := crypto.GeneratePrivateKey()
		egg.SetHash(core.EggplantHasher{})
		sig, err := priKey.Sign(egg.Hash[:])
		assert.Nil(t, err)
		egg.Signature = sig.ToByte()
		egg.FirstSeen = time.Now().Unix()
		go pool.AddEgg(*egg)
	}
	time.Sleep(1 * time.Second)
}

func TestSearchEggByHash(t *testing.T) {
	pool := NewMemoryPool(500)
	data := core.MetaData{EggplantId: 1}
	egg := core.NewEggplant(data)
	hash := egg.SetHash(core.EggplantHasher{})
	go pool.AddEgg(*egg)
	newEgg := pool.SearchEggByHash(hash)
	assert.Equal(t, 1, newEgg.EggplantId)
	assert.Equal(t, *egg, newEgg)

}

func TestIsExist(t *testing.T) {
	pool := NewMemoryPool(500)
	data := core.MetaData{EggplantId: 1}
	egg := core.NewEggplant(data)
	hash := egg.SetHash(core.EggplantHasher{})
	pool.AddEgg(*egg)
	assert.True(t, pool.IsExistEgg(hash))
	assert.False(t, pool.IsExistEgg(types.RandomHash()))
}

func TestDeleteEggByHash(t *testing.T) {
	pool := NewMemoryPool(500)
	data := core.MetaData{EggplantId: 1}
	egg := core.NewEggplant(data)
	hash := egg.SetHash(core.EggplantHasher{})
	pool.AddEgg(*egg)
	assert.True(t, pool.IsExistEgg(hash))
	pool.DeleteEggByHash(hash)
	assert.False(t, pool.IsExistEgg(hash))

}

func TestGoRoutine(t *testing.T) {
	go func() {
		fmt.Printf("go routine time = %v\n", time.Now().UnixNano())
	}()
	fmt.Printf("main thread time = %v\n", time.Now().UnixNano())

}