package args

import (
	"os"
	"strings"
	"sync"
)

var (
	oneArgs     sync.Once
	defaultArgs *Args
)

//Instance desc
//@method Instance desc: command line args instance
//@return (*Args)
func Instance() *Args {
	oneArgs.Do(func() {
		defaultArgs = &Args{make(map[string]interface{})}
	})
	return defaultArgs
}

//Args desc
//@struct Args desc: command line args manager
type Args struct {
	m map[string]interface{}
}

//Parse desc
//@method Parse desc: parse command line args
func (slf *Args) Parse() {
	var tmp []string
	for _, args := range os.Args {
		tmp = append(tmp, args)
	}

	tmp = tmp[1:]

	idx := 0
	for {
		if idx >= len(tmp) {
			break
		}

		cur := tmp[idx]
		idx++
		if strings.HasPrefix(cur, "-") {
			if idx >= len(tmp) {
				slf.m[cur] = true
				break
			}

			next := tmp[idx]
			if strings.HasPrefix(next, "-") {
				continue
			}

			slf.m[cur] = next
			idx++
		}
	}
}

//GetString desc
//@method GetString desc: Return Args Command value
//@param  (string) Command
//@param  (string) default value
//@return (string) value
func (slf *Args) GetString(name string, def string) string {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	return slf.m[name].(string)
}

//GetInt desc
//@method GetInt desc: Return Args Command value
//@param  (string) Command
//@param  (int) default value
//@return (int) value
func (slf *Args) GetInt(name string, def int) int {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	v := slf.m[name]
	if _, ok := v.(int); !ok {
		return def
	}

	return v.(int)
}

//GetInt64 desc
//@method GetInt64 desc: Return Args Command value
//@param  (string) Command
//@param  (int64) default value
//@return (int64) value
func (slf *Args) GetInt64(name string, def int64) int64 {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	v := slf.m[name]
	if _, ok := v.(int64); !ok {
		return def
	}

	return v.(int64)
}

//GetBoolean desc
//@method GetBoolean desc: Return Args Command value
//@param  (string) Command
//@param  (bool) default value
//@return (bool) value
func (slf *Args) GetBoolean(name string, def bool) bool {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	v := slf.m[name]
	if _, ok := v.(bool); !ok {
		return def
	}

	return v.(bool)
}

//GetFloat desc
//@method GetFloat desc: Return Args Command value
//@param  (string) Command
//@param  (float32) default value
//@return (float32) value
func (slf *Args) GetFloat(name string, def float32) float32 {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	v := slf.m[name]
	if _, ok := v.(float32); !ok {
		return def
	}

	return v.(float32)
}

//GetDouble desc
//@method GetDouble desc: Return Args Command value
//@param  (string) Command
//@param  (float64) default value
//@return (float64) value
func (slf *Args) GetDouble(name string, def float64) float64 {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	v := slf.m[name]
	if _, ok := v.(float64); !ok {
		return def
	}

	return v.(float64)
}
