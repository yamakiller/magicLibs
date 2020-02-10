package dbs

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

//RedisDB doc
//@Summary redis db opertioner
//@Struct RedisDB
//@Member (*redis.Pool) a redis connection pool
//@Member (string) redis address
//@Member (int)  redis db code
type RedisDB struct {
	_c    *redis.Pool
	_host string
	_pwd  string
	_db   int
}

//Initial doc
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

	if strings.Index(host, "@") > 0 {
		tmp := strings.Split(host, "@")
		slf._host = tmp[0]
		slf._pwd = tmp[1]
	}

	slf._c = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}

			if slf._pwd != "" {
				if err = c.Send("AUTH", slf._pwd); err != nil {
					c.Close()
					return nil, fmt.Errorf("redis auth error:%+v", err)
				}

			}

			_, err = c.Do("SELECT", db)
			if err != nil {
				c.Close()
				return nil, fmt.Errorf("redis select db error:%+v", err)
			}
			return c, nil
		},
	}

	return nil
}

//Do doc
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

//Close doc
//@Summary close redis db operation
//@Method Close
func (slf *RedisDB) Close() {
	slf._c.Close()
}
