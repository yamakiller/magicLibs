package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/util"
)

func TestPwdValidity(t *testing.T) {
	fmt.Println(util.IsPWDValidity("a12bbbbb"))
}
