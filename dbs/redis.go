package dbs

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

//RedisDB desc
//@Summary redis db opertioner
//@Struct RedisDB
//@Member (*redis.Pool) a redis connection pool
//@Member (string) redis address
//@Member (int)  redis db code
type RedisDB struct {
	_c    *redis.Pool
	_host string
	_db   int
}

//Initial desc
//@Summary initialization Redis DB
//@Method Initial
//@Param (string) redis host address
//@Param (int) redis db code
//@Param (int) redis connection max idle of number
//@Param (int) redis connection max active of number
//@Param (int) redis connection idle time (unit/sec)
//@Return (error) if connecting fail return error ,success return nil
func (slf *RedisDB) Initial(host string, db int, maxIdle int, maxActive int, idleSec int) error {
	slf._host = host
	slf._db = db

	slf._c = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			_, derr := c.Do("SELECT", db)
			if derr != nil {
				panic(fmt.Sprintln("redis select db error:", derr))
			}
			return c, nil
		},
	}

	return nil
}

//Do desc
//@Summary execute redis command
//@Method Do
//@Param (string) command name
//@Param (...interface{}) command params
//@Return (interface{}) execute result
//@Return (error) if execute fail return error, execute success return nil
func (slf *RedisDB) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := slf._c.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

//Close desc
//@Summary close redis db operation
//@Method Close
func (slf *RedisDB) Close() {
	slf._c.Close()
}
