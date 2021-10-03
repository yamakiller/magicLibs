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

//New 创建一个一致性Hash表
func New(replicas int) *Map {
	if replicas <= 0 {
		replicas = 20
	}
	return &Map{_replicas: replicas, _maps: make(map[uint32]interface{})}
}

//ErrEmptyCircle Return an empty
var ErrEmptyCircle = errors.New("empty circle")

type uInt32Slice []uint32

//Len　返回元素个数
func (s uInt32Slice) Len() int {
	return len(s)
}

//Less doc
//@Summary Compare the size of the array i, j position
//@Method Less
//@Param  (int) array index
//@Param  (int) array index
//@Return If the data in the position of the array i is smaller than the data in the j position, it returns True, otherwise it returns False.
func (s uInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

//Swap doc
//@Summary Data exchange between the position of the array i and the data of the j position
//@Method Swap
//@Param  (int) array index
//@Param  (int) array index
func (s uInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Map doc
//@Summary Hash consistency load balancing
//@Struct Map
type Map struct {
	_replicas int
	_keys     uInt32Slice
	_maps     map[uint32]interface{}
	_nodes    []string
	sync.Mutex
}

//IsEmpty 是否空
func (slf *Map) IsEmpty() bool {
	return len(slf._keys) == 0
}

//IsExits 是否存在这个Key
func (slf *Map) IsExits(key string) bool {
	for _, v := range slf._nodes {
		if v == key {
			return true
		}
	}
	return false
}

//Add 增加一个元素
func (slf *Map) Add(key string, v interface{}) {
	for i := 0; i < slf._replicas; i++ {
		slf._maps[getHash([]byte(strconv.Itoa(i)+key))] = v
	}
	idx := slf.getKeyIndex(key)
	if idx == -1 {
		slf._nodes = append(slf._nodes, key)
	} else {
		slf._nodes[idx] = key
	}
	slf.generate()
}

//Remove 删除一个元素
func (slf *Map) Remove(key string) interface{} {
	var v interface{}
	for i := 0; i < slf._replicas; i++ {
		if v == nil {
			v = slf._maps[getHash([]byte(strconv.Itoa(i)+key))]
		}

		delete(slf._maps, getHash([]byte(strconv.Itoa(i)+key)))
	}
	slf.rmKey(key)
	slf.generate()
	return v
}

//Get 获取一个元素[平均返回]
func (slf *Map) Get(name string) (interface{}, error) {
	if len(slf._maps) == 0 {
		return nil, ErrEmptyCircle
	}

	key := getHash([]byte(name))
	i := slf.sreach(key)
	return slf._maps[slf._keys[i]], nil
}

//GetKeys 返回所有的Key
func (slf *Map) GetKeys() []string {
	r := make([]string, len(slf._nodes))
	for k, v := range slf._nodes {
		r[k] = v
	}

	return r
}

func (slf *Map) sreach(key uint32) (i int) {
	i = sort.Search(len(slf._keys), func(x int) bool { return slf._keys[x] >= key })
	if i >= len(slf._keys) {
		i = 0
	}
	return
}

func (slf *Map) getKeyIndex(key string) int {
	for k, v := range slf._nodes {
		if v == key {
			return k
		}
	}
	return -1
}

func (slf *Map) rmKey(key string) {
	for k, v := range slf._nodes {
		if v == key {
			if k == 0 {
				slf._nodes = slf._nodes[1:]
			} else if k == (len(slf._nodes) - 1) {
				slf._nodes = slf._nodes[:k-1]
			} else {
				slf._nodes = append(slf._nodes[:k], slf._nodes[k+1:]...)
			}
			return
		}
	}
}

func (slf *Map) generate() {
	hashes := slf._keys[:0]
	if cap(slf._keys)/(slf._replicas*4) > len(slf._maps) {
		hashes = nil
	}
	for k := range slf._maps {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	slf._keys = hashes
}
