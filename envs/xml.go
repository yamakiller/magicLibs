package envs

import (
	"encoding/xml"

	"github.com/yamakiller/magicLibs/files"
)

//XMLEnv desc
//@struct XMLEnv desc: xml fromat Environment variable loading interface
type XMLEnv struct {
	en map[string]interface{}
}

//Initial desc
//@method Initial desc: Initial json environment variable manager
func (slf *XMLEnv) Initial() {
	slf.en = make(map[string]interface{})
}

//Load desc
//@method Load desc: file load to xml env
//@param  (string) map key
//@param  (string) file name
//@param  (interface{}) env object
//@return (error) load fail return error
func (slf *XMLEnv) Load(key string, fileName string, out interface{}) error {
	fullPathName := files.Instance().GetFullPathForFilename(fileName)
	d, e := files.Instance().GetDataFromFile(fullPathName)
	if e != nil {
		goto fail
	}

	e = xml.Unmarshal(d, out)
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
func (slf *XMLEnv) UnLoad() {
	slf.en = nil
}

//Remove desc
//@method Remove desc: remove environment variable object
//@param  (string) map key
func (slf *XMLEnv) Remove(key string) {
	if _, ok := slf.en[key]; ok {
		delete(slf.en, key)
	}
}

//Get desc
//@method Get desc: Returns environment variable object
//@param  (string) map key
//@return (interface{}) environment variable object
func (slf *XMLEnv) Get(key string) interface{} {
	v, ok := slf.en[key]
	if !ok {
		return nil
	}
	return v
}
