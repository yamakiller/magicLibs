package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/util"
)

func TestUUID(t *testing.T) {
	fmt.Println(util.SpawnUUID())
}