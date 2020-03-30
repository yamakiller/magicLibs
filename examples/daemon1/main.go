package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yamakiller/magicLibs/util"
)

var stdlog, errlog *log.Logger

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func makeFile() {
	// create a simple file (current time).txt
	f, err := os.Create(fmt.Sprintf("%s/%s.txt", os.TempDir(), time.Now().Format(time.RFC3339)))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

type AppTest struct {
}

func (slf *AppTest) Name() string {
	return "app_test"
}

func (slf *AppTest) Desc() string {
	return "app_desc_test"
}

func (slf *AppTest) Open() (string, error) {
	makeFile()
	return "", nil
}

func (slf *AppTest) Close() {
	makeFile()
	stdlog.Println("Got close:")
}

func main() {
	app := &AppTest{}
	d := util.SpawnDaemon(app)
	s, err := d.Open()
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(s)
}
