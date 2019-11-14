package comparator

import "time"

//Comparator desc
//@type Comparator desc: Comparator function
type Comparator func(a, b interface{}) int

//StringComparator desc
//@method StringComparator desc: default a fast comparison on strings
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func StringComparator(a, b interface{}) int {
	s1 := a.(string)
	s2 := b.(string)

	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

//IntComparator desc
//@method IntComparator desc: default a fast comparison on int
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func IntComparator(a, b interface{}) int {
	aAss := a.(int)
	bAss := b.(int)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int8Comparator desc
//@method Int8Comparator desc: default a fast comparison on int8
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Int8Comparator(a, b interface{}) int {
	aAss := a.(int8)
	bAss := b.(int8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int16Comparator desc
//@method Int16Comparator desc: default a fast comparison on int16
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Int16Comparator(a, b interface{}) int {
	aAss := a.(int16)
	bAss := b.(int16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int32Comparator desc
//@method Int32Comparator desc: default a fast comparison on int32
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Int32Comparator(a, b interface{}) int {
	aAss := a.(int32)
	bAss := b.(int32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Int64Comparator desc
//@method Int64Comparator desc: default a fast comparison on int64
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Int64Comparator(a, b interface{}) int {
	aAss := a.(int64)
	bAss := b.(int64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UIntComparator desc
//@method UIntComparator desc: default a fast comparison on uint
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func UIntComparator(a, b interface{}) int {
	aAss := a.(uint)
	bAss := b.(uint)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt8Comparator desc
//@method UInt8Comparator desc: default a fast comparison on uint8
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func UInt8Comparator(a, b interface{}) int {
	aAss := a.(uint8)
	bAss := b.(uint8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt16Comparator desc
//@method UInt16Comparator desc: default a fast comparison on uint16
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func UInt16Comparator(a, b interface{}) int {
	aAss := a.(uint16)
	bAss := b.(uint16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt32Comparator desc
//@method UInt32Comparator desc: default a fast comparison on uint32
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func UInt32Comparator(a, b interface{}) int {
	aAss := a.(uint32)
	bAss := b.(uint32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//UInt64Comparator desc
//@method UInt64Comparator desc: default a fast comparison on uint64
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func UInt64Comparator(a, b interface{}) int {
	aAss := a.(uint64)
	bAss := b.(uint64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Float32Comparator desc
//@method Float32Comparator desc: default a fast comparison on float32
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Float32Comparator(a, b interface{}) int {
	aAss := a.(float32)
	bAss := b.(float32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//Float64Comparator desc
//@method Float64Comparator desc: default a fast comparison on float64
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func Float64Comparator(a, b interface{}) int {
	aAss := a.(float64)
	bAss := b.(float64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//ByteComparator desc
//@method ByteComparator desc: default a fast comparison on byte
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func ByteComparator(a, b interface{}) int {
	aAss := a.(byte)
	bAss := b.(byte)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//RuneComparator desc
//@method RuneComparator desc: default a fast comparison on  time.Time
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func RuneComparator(a, b interface{}) int {
	aAss := a.(rune)
	bAss := b.(rune)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

//TimeComparator desc
//@method TimeComparator desc: default a fast comparison on rune
//@param  (interface{}) a
//@param  (interface{}) b
//@return (int) comparator result
func TimeComparator(a, b interface{}) int {
	aAss := a.(time.Time)
	bAss := b.(time.Time)
	switch {
	case aAss.After(bAss):
		return 1
	case aAss.Before(bAss):
		return -1
	default:
		return 0
	}
}
