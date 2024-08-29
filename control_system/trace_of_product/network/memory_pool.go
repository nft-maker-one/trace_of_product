package network

import (
	"agricultural_meta/core"
	"agricultural_meta/types"
	"agricultural_meta/utils"
	"fmt"
	"sync"
)

// 必须为 fifo 中的哈希索引再增加索引，不然数据量大时 O(N) 的时间复杂度遍历数组然后删除，开销大

// 需要记录数据进入内存池的先后顺序
// 还需要
// 1 数组删除不方便
// 2 golang 自身特性，cap夸大后不会减少
// 3 当通道数据挤压过大时大大减少内存开销
// 4 双向链表结构可以快速添加尾部元素，删除头部元素，适合用作先入先出的通道

type MemoryPool struct {
	Cap       int
	HashIndex map[types.Hash]*ListNode
	fifo      *ListNode
	Eggs      map[types.Hash]core.Eggplant
	lock      sync.RWMutex
}

func NewMemoryPool(cap int) *MemoryPool {

	return &MemoryPool{
		Cap:       cap,
		fifo:      NewHead(),
		Eggs:      make(map[types.Hash]core.Eggplant),
		lock:      sync.RWMutex{},
		HashIndex: make(map[types.Hash]*ListNode),
	}
}

func (m *MemoryPool) AddEgg(egg core.Eggplant) {

	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.Eggs[egg.Hash]; ok {
		utils.LogMsg([]string{"AddEgg"}, []string{fmt.Sprintf("egg [hash:%v] has already in memory pool", egg.Hash)})
		return
	}
	node := NewListNode(egg.Hash)
	m.fifo.TailAdd(node)
	m.HashIndex[egg.Hash] = node
	m.Eggs[egg.Hash] = egg
	utils.LogMsg([]string{"AddEgg"}, []string{fmt.Sprintf("egg [hash:%v] added to memory pool successfully", egg.Hash)})
}

func (m *MemoryPool) PopEgg() (core.Eggplant, error) {
	m.lock.Lock()
	m.lock.Unlock()
	hash, err := m.fifo.TopPop()
	if err != nil {
		utils.LogMsg([]string{"PopEgg"}, []string{"error when deal with message in fifo of memory pool err = " + err.Error()})
		return core.Eggplant{}, err
	}
	egg := m.Eggs[hash]
	delete(m.Eggs, hash)
	delete(m.HashIndex, hash)
	return egg, nil
}

func (m *MemoryPool) SearchEggByHash(hash types.Hash) core.Eggplant {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if _, ok := m.Eggs[hash]; !ok {
		utils.LogMsg([]string{"SearchEggByHash"}, []string{fmt.Sprintf("egg [hash:%v] is not in the memory pool", hash)})
		return core.Eggplant{}
	}
	return m.Eggs[hash]
}

func (m *MemoryPool) IsExistEgg(hash types.Hash) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.Eggs[hash]
	return ok
}

func (m *MemoryPool) DeleteEggByHash(hash types.Hash) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.Eggs[hash]; !ok {
		return
	}
	delete(m.Eggs, hash)
	node := m.HashIndex[hash]
	m.fifo.DeleteNode(node)
	delete(m.HashIndex, hash)
	utils.LogMsg([]string{"DeleteEggByHash"}, []string{fmt.Sprintf("egg[hash:%v] deleted from memory pool successfully", hash)})

}

type ListNode struct {
	prev  *ListNode
	next  *ListNode
	value types.Hash
	len   int
}

func NewHead() *ListNode {
	node := &ListNode{}
	node.prev = node
	node.next = node
	node.value = types.Hash{}
	node.len = 0
	return node
}

func NewListNode(hash types.Hash) *ListNode {
	return &ListNode{
		value: hash,
	}
}

func (n *ListNode) TailAddFromHash(hash types.Hash) {
	tail := n.prev
	newNode := NewListNode(hash)
	tail.next = newNode
	newNode.prev = tail
	n.prev = newNode
	newNode.next = n
	n.len++
}

func (n *ListNode) TailAdd(newNode *ListNode) {
	if newNode == nil {
		utils.LogMsg([]string{"TailAdd"}, []string{"nil node"})
		return
	}
	tail := n.prev
	tail.next = newNode
	newNode.prev = tail
	n.prev = newNode
	newNode.next = n
	n.len++
}

func (n *ListNode) TopPop() (types.Hash, error) {
	if n.len == 0 {
		return types.Hash{}, fmt.Errorf("there is no node in list")
	}
	head := n.next
	second := head.next
	n.next = second
	second.prev = n
	n.len--
	return head.value, nil
}

func (n *ListNode) DeleteNode(node *ListNode) {
	if node == nil {
		return
	}
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
	n.len--
}
