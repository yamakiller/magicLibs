package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/st/hash"
)

func TestHashCon(t *testing.T) {
	h := hash.New(24)
	h.Add("test001", 1)
	h.Add("test002", 2)
	h.Add("test003", 3)
	h.Add("test004", 4)
	h.Add("test005", 5)
	h.Add("test006", 6)

	for i := 0; i < 100; i++ {
		fmt.Println(h.Get(fmt.Sprintf("key%d", i+1)))
	}
}
