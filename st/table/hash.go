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
	_seqID  uint32
	_sz     int
	_arrays []interface{}
}

//Initial desc
//@Method Initial desc: Initialize the hashtable
func (ht *HashTable) Initial() {
	ht._arrays = make([]interface{}, ht.Max)
	ht._seqID = 1
}

//Size desc
//@Method Size desc: Returns the hashtable is number
//@Return (int) size
func (ht *HashTable) Size() int {
	return ht._sz
}

//Push desc
//@Method Push desc: Insert an value
//@Param (interface{}) value
//@Return (int32) insert an value, hash value
//@Return (error)
func (ht *HashTable) Push(v interface{}) (uint32, error) {
	var i uint32
	for i = 0; i < ht.Max; i++ {
		key := ((i + ht._seqID) & ht.Mask)
		hash := key & (ht.Max - 1)
		if ht._arrays[hash] == nil {
			ht._seqID = key + 1
			ht._arrays[hash] = v
			ht._sz++
			return uint32(key), nil
		}
	}

	return 0, ErrHashTableFulled
}

//Get desc
//@Method Get desc: Returns the one elements from the hashtable
//@Param  (uint32) hash value
//@Return (interface{})
func (ht *HashTable) Get(key uint32) interface{} {
	hash := key & uint32(ht.Max-1)
	if ht._arrays[hash] != nil && ht.Comp(ht._arrays[hash], key) == 0 {
		return ht._arrays[hash]
	}
	return nil
}

//GetValues desc
//@Method GetValues desc: Returns the elements of all from hashtable
//@Return ([]interface{}) Returns all value
func (ht *HashTable) GetValues() []interface{} {
	if ht._sz == 0 {
		return nil
	}

	i := 0
	result := make([]interface{}, ht._sz)
	for _, v := range ht._arrays {
		if v == nil {
			continue
		}
		result[i] = ht._arrays[i]
		i++
	}

	return result
}

//Remove desc
//@Method Remove desc: removes one elements in the hashtable
//@Param  (uint32) hash value
//@Return (bool)
func (ht *HashTable) Remove(key uint32) bool {
	hash := uint32(key) & uint32(ht.Max-1)
	if ht._arrays[hash] != nil && ht.Comp(ht._arrays[hash], key) == 0 {
		ht._arrays[hash] = nil
		ht._sz--
		return true
	}

	return false
}
