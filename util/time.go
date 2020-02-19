package util

import "time"

//TimeNowFormat Returns　当前时间字符串
//@Summary get now time and format
//@Return (string) YYYY-MM-DD hh:ii:ss
func TimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//Timestamp Returns 当前时间毫秒
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
