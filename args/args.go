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
//@Method Instance desc: command line args instance
//@Return (*Args)
func Instance() *Args {
	oneArgs.Do(func() {
		defaultArgs = &Args{make(map[string]interface{})}
	})
	return defaultArgs
}

//Args desc
//@Struct Args desc: command line args manager
type Args struct {
	m map[string]interface{}
}

//Parse desc
//@Method Parse desc: parse command line args
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
//@Method GetString desc: Return Args Command value
//@Param  (string) Command
//@Param  (string) default value
//@Return (string) value
func (slf *Args) GetString(name string, def string) string {
	if _, ok := slf.m[name]; !ok {
		return def
	}

	return slf.m[name].(string)
}

//GetInt desc
//@Method GetInt desc: Return Args Command value
//@Param  (string) Command
//@Param  (int) default value
//@Return (int) value
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
//@Method GetInt64 desc: Return Args Command value
//@Param  (string) Command
//@Param  (int64) default value
//@Return (int64) value
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
//@Method GetBoolean desc: Return Args Command value
//@Param  (string) Command
//@Param  (bool) default value
//@Return (bool) value
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
//@Method GetFloat desc: Return Args Command value
//@Param  (string) Command
//@Param  (float32) default value
//@Return (float32) value
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
//@Method GetDouble desc: Return Args Command value
//@Param  (string) Command
//@Param  (float64) default value
//@Return (float64) value
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
