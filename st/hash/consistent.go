package hash

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

func getHash(str []byte) uint32 {
	p := uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(str); i++ {
		hash = (hash ^ uint32(str[i])) * p
	}
	hash += hash << 13
	hash ^= hash >> 7
	hash += hash << 3
	hash ^= hash >> 17
	hash += hash << 5
	return hash
}

//NewConsistentHash doc
//@Method NewConsistentHash @Summary Create a hash consistent loader
//@Param  (int) replicas of number
//@Return (*Map)
func NewConsistentHash(replicas int) *Map {
	if replicas <= 0 {
		replicas = 20
	}
	return &Map{_replicas: replicas,
		_circle: make(map[uint32]interface{})}
}

//ErrEmptyCircle Return an empty
var ErrEmptyCircle = errors.New("empty circle")

type uInt32Slice []uint32

//Len doc
//@Method Len @Summary array lenght
//@Return (int)
func (s uInt32Slice) Len() int {
	return len(s)
}

//Less doc
//@Method Less @Summary Compare the size of the array i, j position
//@Param  (int) array index
//@Param  (int) array index
//@Return If the data in the position of the array i is smaller than the data in the j position, it returns True, otherwise it returns False.
func (s uInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

//Swap doc
//@Method Swap @Summary Data exchange between the position of the array i and the data of the j position
//@Param  (int) array index
//@Param  (int) array index
func (s uInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Map doc
//@Struct Map @Summary Hash consistency load balancing
type Map struct {
	_replicas     int
	_sortedHashes uInt32Slice
	_circle       map[uint32]interface{}
	sync.RWMutex
}

//UnAdd doc
//@Method UnAdd @Summary Join an object, not locked
//@Param (string) key
//@Param (interface{}) element value
func (m *Map) UnAdd(key string, v interface{}) {
	for i := 0; i < m._replicas; i++ {
		m._circle[getHash([]byte(strconv.Itoa(i)+key))] = v
	}
	m.updateSortedHashes()
}

//UnRemove doc
//@Method UnRemove @Summary Delete an object, not locked
//@Param (string) key
func (m *Map) UnRemove(key string) {
	for i := 0; i < m._replicas; i++ {
		delete(m._circle, getHash([]byte(strconv.Itoa(i)+key)))
	}
	m.updateSortedHashes()
}

//UnGet doc
//@Method UnGet @Summary Return an object, not locked
//@Param  (string) name
//@Return (interface{}) element value
//@Return (error)
func (m *Map) UnGet(name string) (interface{}, error) {
	if len(m._circle) == 0 {
		return "", ErrEmptyCircle
	}

	key := getHash([]byte(name))
	i := m.sreach(key)
	return m._circle[m._sortedHashes[i]], nil
}

func (m *Map) sreach(key uint32) (i int) {
	f := func(x int) bool {
		return m._sortedHashes[x] > key
	}
	i = sort.Search(len(m._sortedHashes), f)
	if i >= len(m._sortedHashes) {
		i = 0
	}
	return
}

func (m *Map) updateSortedHashes() {
	hashes := m._sortedHashes[:0]
	if cap(m._sortedHashes)/(m._replicas*4) > len(m._circle) {
		hashes = nil
	}
	for k := range m._circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	m._sortedHashes = hashes
}
