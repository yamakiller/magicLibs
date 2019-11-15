package table

import (
	"errors"

	"github.com/yamakiller/magicLibs/st/comparator"
)

var (
	//ErrHashTableFulled Table is full
	ErrHashTableFulled = errors.New("Table is full")
)

//HashTable Hash allocation table
type HashTable struct {
	Mask   uint32
	Max    uint32
	Comp   comparator.Comparator
	seqID  uint32
	sz     int
	arrays []interface{}
}

//Initial desc
//@method Initial desc: Initialize the hashtable
func (ht *HashTable) Initial() {
	ht.arrays = make([]interface{}, ht.Max)
	ht.seqID = 1
}

//Size desc
//@method Size desc: Returns the hashtable is number
//@return (int) size
func (ht *HashTable) Size() int {
	return ht.sz
}

//Push desc
//@method Push desc: Insert an value
//@param (interface{}) value
//@return (int32) insert an value, hash value
//@return (error)
func (ht *HashTable) Push(v interface{}) (uint32, error) {
	var i uint32
	for i = 0; i < ht.Max; i++ {
		key := ((i + ht.seqID) & ht.Mask)
		hash := key & (ht.Max - 1)
		if ht.arrays[hash] == nil {
			ht.seqID = key + 1
			ht.arrays[hash] = v
			ht.sz++
			return uint32(key), nil
		}
	}

	return 0, ErrHashTableFulled
}

//Get desc
//@method Get desc: Returns the one elements from the hashtable
//@param  (uint32) hash value
//@return (interface{})
func (ht *HashTable) Get(key uint32) interface{} {
	hash := key & uint32(ht.Max-1)
	if ht.arrays[hash] != nil && ht.Comp(ht.arrays[hash], key) == 0 {
		return ht.arrays[hash]
	}
	return nil
}

//GetValues desc
//@method GetValues desc: Returns the elements of all from hashtable
//@return ([]interface{}) Returns all value
func (ht *HashTable) GetValues() []interface{} {
	if ht.sz == 0 {
		return nil
	}

	i := 0
	result := make([]interface{}, ht.sz)
	for _, v := range ht.arrays {
		if v == nil {
			continue
		}
		result[i] = ht.arrays[i]
		i++
	}

	return result
}

//Remove desc
//@method Remove desc: removes one elements in the hashtable
//@param  (uint32) hash value
//@return (bool)
func (ht *HashTable) Remove(key uint32) bool {
	hash := uint32(key) & uint32(ht.Max-1)
	if ht.arrays[hash] != nil && ht.Comp(ht.arrays[hash], key) == 0 {
		ht.arrays[hash] = nil
		ht.sz--
		return true
	}

	return false
}
