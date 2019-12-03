package util

import "time"

//TimeNowFormat doc
//@Method TimeNowFormat @Summary get now time and format
//@Return (string) YYYY-MM-DD hh:ii:ss
func TimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
