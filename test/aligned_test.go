package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/mmath"
)

func TestAligned(t *testing.T) {
	fmt.Println(mmath.Aligned(1400))
	fmt.Println(mmath.Align(1401, 4))
}
