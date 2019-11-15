package envs

import (
	"github.com/yamakiller/magicLibs/files"
	"github.com/yamakiller/magicLibs/util"
)

//JSONEnv desc
//@struct JSONEnv desc: json fromat Environment variable loading interface
//
type JSONEnv struct {
	en map[string]interface{}
}

//Initial desc
//@method Initial desc: Initial json environment variable manager
func (slf *JSONEnv) Initial() {
	slf.en = make(map[string]interface{})
}

//Load desc
//@method Load desc: file load to json env
//@param  (string) map key
//@param  (string) file name
//@param  (interface{}) env object
//@return (error) load fail return error
func (slf *JSONEnv) Load(key string, fileName string, out interface{}) error {
	fullPathName := files.Instance().GetFullPathForFilename(fileName)
	d, e := files.Instance().GetDataFromFile(fullPathName)
	if e != nil {
		goto fail
	}

	e = util.JSONUnSerialize(d, out)
	if e != nil {
		goto fail
	}

	slf.en[key] = out
	return nil
fail:
	slf.en[key] = out
	return e
}

//UnLoad desc
//@method UnLoad desc:
func (slf *JSONEnv) UnLoad() {
	slf.en = nil
}

//Remove desc
//@method Remove desc: remove environment variable object
//@param  (string) map key
func (slf *JSONEnv) Remove(key string) {
	if _, ok := slf.en[key]; ok {
		delete(slf.en, key)
	}
}

//Get desc
//@method Get desc: Returns environment variable object
//@param  (string) map key
//@return (interface{}) environment variable object
func (slf *JSONEnv) Get(key string) interface{} {
	v, ok := slf.en[key]
	if !ok {
		return nil
	}
	return v
}
