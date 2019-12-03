package util

//IsPower doc
//@Method IsPower @Summary Determine if the value is a power of 2
//@Param  (int)  check value
//@Return (bool) yes:true no:false
func IsPower(n int) bool {
	if n < 2 {
		return false
	}

	if (n&n - 1) == 0 {
		return true
	}
	return false
}
