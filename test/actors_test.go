package test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/log"
)

type TActor struct {
}

func (slf *TActor) Receive(ctx *actors.Context) {
	ctx.Info("Log Receive")
}

func TestActorsClose(t *testing.T) {
	//
	hlog := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "white+h",
		TimestampStyle: "black+h"})
	hlog.SetFormatter(formatter)
	hlog.SetOutput(os.Stdout)

	logSystem := &log.DefaultAgent{}
	logSystem.WithHandle(hlog)

	engine := actors.New(nil)
	engine.WithLogger(logSystem)

	PID1, err := engine.New(func(pid *actors.PID) actors.Actor {
		return &TActor{}
	})
	if err != nil {
		logSystem.Info("", "创建Actor错误")
	}

	PID1.Stop()

	logSystem.Info("", "测试1")
	engine.Close()
	logSystem.Close()

}

type NewTest struct {
	_engine *actors.Core
	_gw     sync.WaitGroup
}

func (slf *NewTest) run() {
	defer slf._gw.Done()
	for i := 0; i < 100; i++ {
		PID1, err := slf._engine.New(func(pid *actors.PID) actors.Actor {
			return &TActor{}
		})
		if err != nil {
			fmt.Printf("创建失败:%+v", err)
		}

		PID1.Stop()
	}
}

func TestActorsSyncNew(t *testing.T) {
	hlog := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "white+h",
		TimestampStyle: "black+h"})
	hlog.SetFormatter(formatter)
	hlog.SetOutput(os.Stdout)

	logSystem := &log.DefaultAgent{}
	logSystem.WithHandle(hlog)

	engine := actors.New(nil)
	engine.WithLogger(logSystem)

	testA := &NewTest{_engine: engine}
	testA._gw.Add(2)
	go testA.run()
	go testA.run()
	testA._gw.Wait()

	engine.Close()
	logSystem.Close()
}
