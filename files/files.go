package magicFiles

import (
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
)

type Files struct {
	_dir        Directory
	_cached     sync.Map
	_cachedFils int64
	_cacheMem   int64
}

func (slf *Files) Close() {

	var key []string
	slf._cached.Range(func(k, v interface{}) bool {
		key = append(key, k.(string))
		return true
	})

	for _, k := range key {
		slf._cached.Delete(k)
	}
	slf._cachedFils = 0
	slf._cacheMem = 0
	key = nil
}

//IsFileExist desc
//@method IsFileExist desc: files is exist
//@param (string) file path and name
//@return (bool) exist:true not exits: false
func (slf *Files) IsFileExist(fullPath string) bool {
	h := slf.getCache(fullPath)
	if h != nil {
		return true
	}

	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return true
}

//DeleteFile desc
//@method DeleteFile desc: delete file and remove cache
//@param  (string) full path and file name
//@retur  (error) delete fail error
func (slf *Files) DeleteFile(fullPath string) error {
	h := slf.getCache(fullPath)
	if h != nil {
		atomic.AddInt64(&slf._cacheMem, -int64(h._size))
		atomic.AddInt64(&slf._cachedFils, -1)
		slf._cached.Delete(fullPath)
	}
	err := os.Remove(fullPath)
	if err != nil {
		return err
	}
	return nil
}

//GetCacheFiles desc
//@method GetCacheFiles desc: Return Cache file of number
//@return (int) file of number
func (slf *Files) GetCacheFiles() int {
	return int(atomic.LoadInt64(&slf._cachedFils))
}

//GetCacheMem desc
//@method GetCacheMem desc: Return Cache memory
//@return (int64) Cache size
func (slf *Files) GetCacheMem() int64 {
	return atomic.LoadInt64(&slf._cacheMem)
}

//GetDataFromFile desc
//@method GetDataFromFile desc: Use Cache or Read File Return data
//@param  (string) full path and file name
//@return (*FileHandle) file handle
//@return (error)
func (slf *Files) GetDataFromFile(fullPath string) (*FileHandle, error) {
	h := slf.getCache(fullPath)
	if h != nil {
		return h, nil
	}

	d, e := ioutil.ReadFile(fullPath)
	if e != nil {
		return nil, e
	}

	h = &FileHandle{d, len(d)}
	oh, ok := slf._cached.LoadOrStore(fullPath, h)
	if ok {
		return oh.(*FileHandle), nil
	}
	atomic.AddInt64(&slf._cacheMem, int64(len(d)))
	atomic.AddInt64(&slf._cachedFils, 1)
	return h, nil
}

//getCache desc
//@method getCache desc: Return file cache
//@return (*FileHandle) file cache handle
func (slf *Files) getCache(fullPath string) *FileHandle {
	r, ok := slf._cached.Load(fullPath)
	if !ok {
		return nil
	}

	h, ok := r.(*FileHandle)
	if !ok {
		return nil
	}
	return h
}
