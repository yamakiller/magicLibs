package files

import (
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
)

//Files desc
//@struct Files desc: files system
type Files struct {
	_dir        Directory
	_cached     sync.Map
	_cachedFils int64
	_cacheMem   int64
}

//Initial desc
//@method Initial desc: Initial files system
func (slf *Files) Initial() {
	slf._dir.Initial()
}

//Close desc
//@method Close desc: Close System and clear data
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

//GetRoot desc
//@method GetRoot desc: Return root dir
//@return (string) dir
func (slf *Files) GetRoot() string {
	return slf._dir.rootPath
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
//@method GetDataFromFile desc: Retrun File data
//@param  (string) full path and file name
//@return ([]byte) file data
//@return (error)
func (slf *Files) GetDataFromFile(fullPath string) ([]byte, error) {

	d, e := ioutil.ReadFile(fullPath)
	if e != nil {
		return nil, e
	}

	return d, nil
}

//GetDataFromCacheFile desc
//@method GetDataFromCacheFile desc: Use Cache or Read File Return data
//@param  (string) full path and file name
//@return (*FileHandle) file handle
//@return (error)
func (slf *Files) GetDataFromCacheFile(fullPath string) (*FileHandle, error) {
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

//WithRoot desc
//@method WithRoot desc:setting system root dir
//@param (string) path
func (slf *Files) WithRoot(path string) {
	slf._dir.WithRoot(path)
}

//GetFullPathForFilename desc
//@method GetFullPathForFilename desc: Returns file full path and file name
//@param  (string) file path
//@return (string) full path and file name
func (slf *Files) GetFullPathForFilename(filePath string) string {
	return slf._dir.GetFullPathName(filePath)
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

var (
	oneFiles    sync.Once
	defaultFile *Files
)

//Instance desc
//@method Instance desc: Files System instance
//@return (*Files) Files System object
func Instance() *Files {
	oneFiles.Do(func() {
		defaultFile = &Files{}
		defaultFile.Initial()
	})

	return defaultFile
}
