package hash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uintArray []uint32

//Len desc
//@method Len desc: array lenght
//@return (int)
func (x uintArray) Len() int {
	return len(x)
}

//Less desc
//@method Less desc: Compare the size of the array i, j position
//@param  (int) array index
//@param  (int) array index
//@return If the data in the position of the array i is smaller than the data in the j position, it returns True, otherwise it returns False.
func (x uintArray) Less(i, j int) bool {
	return x[i] < x[j]
}

//Swap desc
//@method Swap desc: Data exchange between the position of the array i and the data of the j position
//@param  (int) array index
//@param  (int) array index
func (x uintArray) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//Consistent desc
//@struct Consistent desc: Hash consistency
type Consistent struct {
	NumberOfReplicas int
	sortedHashes     uintArray
	//
	circle map[uint32]interface{}
	size   int
	sync.RWMutex
}

//NewConsistent desc
//@method NewConsistent desc: create Consistent object
//@param  (int) replicas of number
func NewConsistent(n int) *Consistent {
	return &Consistent{NumberOfReplicas: n, circle: make(map[uint32]interface{})}
}

//Push desc
//@method Push desc: inserts a sring element in the consistent hash.
//@param  (string) element key
//@param  (interface{}) element object
func (slf *Consistent) Push(e string, v interface{}) {
	slf.Lock()
	defer slf.Unlock()
	slf.push(e, v)
}

//Erase desc
//@method Erase desc: removes an element in the consistent hash.
//@param (string) element key
func (slf *Consistent) Erase(e string) {
	slf.Lock()
	defer slf.Unlock()
	slf.erase(e)
}

//Get desc
//@method Get desc: returns an element close to where name hashes to in the circle
//@param  (string) element key/name
//@return (interface{}) element object
//@return (error)
func (slf *Consistent) Get(name string) (interface{}, error) {
	slf.RLock()
	defer slf.RUnlock()

	if len(slf.circle) == 0 {
		return nil, errors.New("empty consistent")
	}
	key := slf.hashCalc(name)
	i := slf.search(key)
	return slf.circle[slf.sortedHashes[i]], nil
}

//Sreach desc
//@method Sreach desc: return an element to where f returns 0 to in the circle
//@param (interface{}) key
//@param (func(key interface{}, val interface{}) int) function
func (slf *Consistent) Sreach(key interface{}, f func(key interface{}, val interface{}) int) interface{} {
	slf.RLock()
	defer slf.RUnlock()

	for _, v := range slf.circle {
		if f(key, v) == 0 {
			return v
		}
	}
	return nil
}

//Range desc
//@method Range desc: Traverse access to all elements
//@param  (func(val interface{})) call function
func (slf *Consistent) Range(f func(val interface{})) {
	slf.RLock()
	defer slf.RUnlock()

	for _, v := range slf.circle {
		f(v)
	}
}

//Size desc
//@method Size desc: Returns memory number to in the circle
//@return (int) number
func (slf *Consistent) Size() int {
	return slf.size
}

func (slf *Consistent) push(e string, v interface{}) {
	for i := 0; i < slf.NumberOfReplicas; i++ {
		slf.circle[slf.hashCalc(slf.genKey(e, i))] = v
	}

	slf.updateSortedHashes()
	slf.size++
}

func (slf *Consistent) erase(e string) {
	for i := 0; i < slf.NumberOfReplicas; i++ {
		delete(slf.circle, slf.hashCalc(slf.genKey(e, i)))
	}
	slf.updateSortedHashes()
	slf.size--
}

func (slf *Consistent) search(key uint32) (i int) {
	f := func(x int) bool {
		return slf.sortedHashes[x] > key
	}

	i = sort.Search(len(slf.sortedHashes), f)
	if i >= len(slf.sortedHashes) {
		i = 0
	}
	return i
}

// generates a string key for an element with an index
func (slf *Consistent) genKey(s string, idx int) string {
	return strconv.Itoa(idx) + s
}

func (slf *Consistent) hashCalc(s string) uint32 {
	if len(s) < 64 {
		var scratch [64]byte
		copy(scratch[:], s)
		return crc32.ChecksumIEEE(scratch[:len(s)])
	}
	return crc32.ChecksumIEEE([]byte(s))
}

func (slf *Consistent) updateSortedHashes() {
	hashes := slf.sortedHashes[:0]
	if cap(slf.sortedHashes)/(slf.NumberOfReplicas*4) > len(slf.circle) {
		hashes = nil
	}
	for k := range slf.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	slf.sortedHashes = hashes
}
