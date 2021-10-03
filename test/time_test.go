package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/util"
)

func TestTime(t *testing.T) {
	prvDay := util.ToDate("2020-02-18")
	nxtDay := util.ToDate("2020-02-20")

	fmt.Println("当前时间:", util.TimestampDay())
	fmt.Println("昨天时间:", util.ToDay(prvDay))
	fmt.Println("明天时间:", util.ToDay(nxtDay))
}
