package util

//IsPower desc
//@method IsPower desc: Determine if the value is a power of 2
//@param  (int)  check value
//@return (bool) yes:true no:false
func IsPower(n int) bool {
	if n < 2 {
		return false
	}

	if (n&n - 1) == 0 {
		return true
	}
	return false
}
