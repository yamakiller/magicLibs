package util

import (
	"math/rand"
	"regexp"
)

var _letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//SubStr desc
//@Method SubStr desc:Cut string ends with length
//@Param  (string) source string
//@Param  (int)    start pos
//@Param  (int)    sub length
//@Return (string) sub string
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
//@Method SubStr2 desc : Cut string ends with index
//@Param  (string) source string
//@Param  (int)    start pos
//@Param  (int)    end pos
//@Return (string) sub string
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

//RandString desc
//@Method RandStr desc : Randomly generate a string of length n
//@Param (int) length
//@Return (string)
func RandStr(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = _letterRunes[rand.Intn(len(_letterRunes))]
	}
	return string(b)
}

//IsPWDValidity desc
//@Method IsPWDValidity desc: Verify password is valid
//@Param (string) password
//@Return (bool) is valid
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

//IsAccountValidity desc
//@Method IsACCValidity desc: Verify account is valid
//@Param (string) account
//@Return (bool) is valid
func IsAccountValidity(account string) bool {
	b, e := regexp.MatchString("^[a-zA-Z0-9_-]{8,16}$", account)
	if e != nil {
		panic(e)
	}
	return b
}

//IsCaptchaValidity desc
//@Method IsCaptchaValidity desc: Verify captcha is valid
//@Param (string) captcha
//@Return (bool) is valid
func IsCaptchaValidity(captcha string) bool {
	b, e := regexp.MatchString("^[0-9]{6,6}", captcha)
	if e != nil {
		panic(e)
	}
	return b
}
