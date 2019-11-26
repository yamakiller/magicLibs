package util

import (
	"regexp"
)

//SubStr desc
//@method SubStr desc:Cut string ends with length
//@param  (string) source string
//@param  (int)    start pos
//@param  (int)    sub length
//@return (string) sub string
func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//SubStr2 desc
//@method SubStr2 desc : Cut string ends with index
//@param  (string) source string
//@param  (int)    start pos
//@param  (int)    end pos
//@return (string) sub string
func SubStr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}

//IsPWDValidity desc
//@method IsPWDValidity desc: Verify password is valid
//@param (string) password
//@return (bool) is valid
func IsPWDValidity(pwd string) bool {
	b, e := regexp.MatchString("^([a-zA-Z_-].*)([0-9].*)$", pwd)
	if e != nil {
		panic(e)
	}

	if b {
		if len(pwd) >= 8 && len(pwd) <= 16 {
			return true
		}
		return false
	}

	return false
}

//IsACCValidity desc
//@method IsACCValidity desc: Verify account is valid
//@param (string) account
//@return (bool) is valid
func IsACCValidity(account string) bool {
	b, e := regexp.MatchString("^[a-zA-Z0-9_-]{8,16}$", account)
	if e != nil {
		panic(e)
	}
	return b
}
