package mutex

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

//RedisMethodDo desc
//@type (func(commandName string, args ...interface{}) (interface{}, error))
type RedisMethodDo func(commandName string, args ...interface{}) (interface{}, error)

//RedisMutex desc
//@struct RedisMutex desc: redis mutex object
//@member (string) resource
//@member (string) token
//@member (int)    timeout
type RedisMutex struct {
	_doFun    RedisMethodDo
	_resource string
	_token    string
	_timeout  int
}

func (slf *RedisMutex) tryLock() (ok bool, err error) {
	_, err = slf._doFun("SET", slf.key(), slf._token, "EX", int(slf._timeout), "NX")
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (slf *RedisMutex) key() string {
	return fmt.Sprintf("redislock:%s", slf._resource)
}

//Unlock desc
//@method Unlock desc: unlocking
//@return (error) unlock fail returns error informat
func (slf *RedisMutex) Unlock() (err error) {
	_, err = slf._doFun("del", slf.key())
	return
}

//AddTimeOut desc
//@method AddTimeOut desc: rest/append lock timeout time
//@param (int64) setting/append time
//@return (bool)
//@reutrn (error)
func (slf *RedisMutex) AddTimeOut(exTime int64) (ok bool, err error) {
	ttl, err := redis.Int64(slf._doFun("TTL", slf.key()))
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl > 0 {
		_, err := redis.String(slf._doFun("SET", slf.key(), slf._token, "EX", int(ttl+exTime)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

//TryRedisLock Try to acquire a lock
//@method TryRedisLock desc: Try to acquire a locking
//@param (RedisMethodDo) redis do function
//@param (string) lock object(key/name)
//@param (string) lock token
//@param (int)    lock timeout millsec
//@return (*RedisMutex) redis mutex object
//@return (bool) redis lock is success
//@return (error) redis lock fail error informat
func TryRedisLock(doFun RedisMethodDo,
	resouse string,
	token string,
	timeout int) (m *RedisMutex, ok bool, err error) {
	return TryRedisLockWithTimeOut(doFun, resouse, token, timeout)
}

//TryRedisLockWithTimeOut desc
//@method TryRedisLockWithTimeOut desc: Try to acquire the lock and set the lock timeout
//@param (RedisMethodDo) redis do function
//@param (string) lock object(key/name)
//@param (string) lock token
//@param (int)    lock timeout millsec
//@return (*RedisMutex) redis mutex object
//@return (bool) redis lock is success
//@return (error) redis lock fail error informat
func TryRedisLockWithTimeOut(doFun RedisMethodDo,
	resouse string,
	token string,
	timeout int) (m *RedisMutex, ok bool, err error) {
	m = &RedisMutex{doFun, resouse, token, timeout}
	ok, err = m.tryLock()
	if !ok || err != nil {
		m = nil
	}
	return
}
