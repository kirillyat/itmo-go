package consistenthash

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"
)

// Node represents a node in the consistent hash ring.
type Node interface {
	ID() string
}

// ConsistentHash manages the consistent hashing ring.
type ConsistentHash[N Node] struct {
	sync.RWMutex
	hashRing   []int      // Sorted list of hash values representing virtual nodes
	nodeMap    map[int]*N // Map from hash value to actual node
	virtualRep int        // Number of virtual nodes per actual node
}

// New creates a new ConsistentHash instance with a default number of virtual nodes per node.
func New[N Node]() *ConsistentHash[N] {
	return &ConsistentHash[N]{
		hashRing:   []int{},
		nodeMap:    make(map[int]*N),
		virtualRep: 100, // default number of virtual nodes per real node
	}
}

// AddNode adds a node to the consistent hash ring.
func (h *ConsistentHash[N]) AddNode(n *N) {
	h.Lock()
	defer h.Unlock()
	for i := 0; i < h.virtualRep; i++ {
		hash := h.hashKey(fmt.Sprintf("%s#%d", (*n).ID(), i))
		h.hashRing = insertSorted(h.hashRing, hash)
		h.nodeMap[hash] = n
	}
}

// RemoveNode removes a node from the consistent hash ring.
func (h *ConsistentHash[N]) RemoveNode(n *N) {
	h.Lock()
	defer h.Unlock()
	for i := 0; i < h.virtualRep; i++ {
		hash := h.hashKey(fmt.Sprintf("%s#%d", (*n).ID(), i))
		h.hashRing = removeSorted(h.hashRing, hash)
		delete(h.nodeMap, hash)
	}
}

// GetNode gets the appropriate node for the given key.
func (h *ConsistentHash[N]) GetNode(key string) *N {
	h.RLock()
	defer h.RUnlock()
	if len(h.hashRing) == 0 {
		return nil
	}
	hash := h.hashKey(key)
	index := findNearestHash(h.hashRing, hash)
	return h.nodeMap[h.hashRing[index]]
}

// hashKey computes the hash for a given key using SHA-256.
func (h *ConsistentHash[N]) hashKey(key string) int {
	hash := sha256.Sum256([]byte(key))
	return int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])
}

// insertSorted inserts a value into a sorted slice and returns the new sorted slice.
func insertSorted(slice []int, value int) []int {
	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
	slice = append(slice, 0)
	copy(slice[i+1:], slice[i:])
	slice[i] = value
	return slice
}

// removeSorted removes a value from a sorted slice and returns the new sorted slice.
func removeSorted(slice []int, value int) []int {
	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
	if i < len(slice) && slice[i] == value {
		slice = append(slice[:i], slice[i+1:]...)
	}
	return slice
}

// findNearestHash finds the nearest hash in the sorted slice.
func findNearestHash(slice []int, hash int) int {
	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= hash })
	if i == len(slice) {
		return 0
	}
	return i
}
