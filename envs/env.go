package envs

import "sync"

var (
	oneEnv     sync.Once
	defaultEnv IEnv
)

//With desc
//@method With desc: Setting Environment variable object
//@param (IEnv) new Environment variable object
func With(h IEnv) {
	defaultEnv = h
}

//Instance desc
//@method Instance desc: Environment variable object instance
//@return (IEnv) an Environment variable object
func Instance() IEnv {
	oneEnv.Do(func() {
		if defaultEnv == nil {
			defaultEnv = &JSONEnv{}
		}
		defaultEnv.Initial()
	})

	return defaultEnv
}

//IEnv desc
//@interface IEnv desc: Environment variable loading interface
//@method Load desc: file load to environment variable
type IEnv interface {
	Initial()
	Load(key string, fileName string, out interface{}) error
	UnLoad()
	Remove(key string)
	Get(key string) interface{}
}
