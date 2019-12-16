package mmath

//Aligned doc
//@Summary aligned byte
func Aligned(n int) int {

	if (n & (n - 1)) != 0 {
		n = n | (n >> 1)
		n = n | (n >> 2)
		n = n | (n >> 4)
		n = n | (n >> 8)
		n = n | (n >> 16)
		n++
	}

	return n
}
