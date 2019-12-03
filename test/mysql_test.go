package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/dbs"
)

//TestMySQL doc
func TestMySQL(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test"
	hsql := dbs.MySQLDB{}
	dep := dbs.MySQLDeploy{DSN: dsn, Max: 1, Idle: 1, Life: 1000 * 60 * 10}
	if err := dbs.DoMySQLDeploy(&hsql, &dep); err != nil {
		fmt.Println("连接失败:", err)
		return
	}

	result, err := hsql.Query("select c1,c2,c3,c4,c5,c6 from test_1 where c1=?", 1)
	if err != nil {
		fmt.Println("查询失败:", err)
		hsql.Close()
		return
	}

	for result.Next() {
		a, _ := result.GetValue(0)
		b, _ := result.GetValue(1)
		c, _ := result.GetValue(2)
		d, _ := result.GetValue(3)
		e, _ := result.GetValue(4)
		f, _ := result.GetValue(5)
		a.Print()
		b.Print()
		c.Print()
		d.Print()
		e.Print()
		f.Print()
	}

	result.Close()
	hsql.Close()
}
