package dbs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

//RedisOptions redis配置参数
type RedisOptions struct {
	Network     string                                 // 通讯协议，默认为 tcp
	Addr        string                                 // redis服务的地址，默认为 127.0.0.1:6379
	Password    string                                 // redis鉴权密码
	Db          int                                    // 数据库
	MaxActive   int                                    // 最大活动连接数，值为0时表示不限制
	MaxIdle     int                                    // 最大空闲连接数
	IdleTimeout int                                    // 空闲连接的超时时间，超过该时间则关闭连接。单位为秒。默认值是5分钟。值为0时表示不关闭空闲连接。此值应该总是大于redis服务的超时时间。
	Prefix      string                                 // 键名前缀
	Marshal     func(v interface{}) ([]byte, error)    // 数据序列化方法，默认使用json.Marshal序列化
	Unmarshal   func(data []byte, v interface{}) error // 数据反序列化方法，默认使用json.Unmarshal序列化
}

//RedisDB doc
//@Summary redis db opertioner
//@Struct RedisDB
//@Member (*redis.Pool) a redis connection pool
//@Member (string) redis address
//@Member (int)  redis db code
type RedisDB struct {
	_c         *redis.Pool
	_host      string
	_pwd       string
	_db        int
	_prefix    string
	_marshal   func(v interface{}) ([]byte, error)
	_unmarshal func(data []byte, v interface{}) error
}

//Initial doc
//@Summary initialization Redis DB
//@Method Initial
//@Param (string) redis host address
//@Param (string) redis password
//@Param (int) redis db code
//@Param (int) redis connection max idle of number
//@Param (int) redis connection max active of number
//@Param (int) redis connection idle time (unit/sec)
//@Return (error) if connecting fail return error ,success return nil
func (slf *RedisDB) Initial(host, pwd string, db int, maxIdle int, maxActive int, idleSec int) error {
	slf._host = host
	slf._db = db
	slf._pwd = pwd
	slf._marshal = json.Marshal
	slf._unmarshal = json.Unmarshal

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
			c, err := redis.Dial("tcp", slf._host)
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

//Initial2 doc
//@Summary initialization Redis DB
//@Method Initial
//@Return (error) if connecting fail return error ,success return nil
func (slf *RedisDB) Initial2(options interface{}) error {
	switch opts := options.(type) {
	case RedisOptions:
		if opts.Network == "" {
			opts.Network = "tcp"
		}
		if opts.Addr == "" {
			opts.Addr = "127.0.0.1:6379"
		}
		if opts.MaxIdle == 0 {
			opts.MaxIdle = 3
		}
		if opts.IdleTimeout == 0 {
			opts.IdleTimeout = 300
		}
		if opts.Marshal == nil {
			slf._marshal = json.Marshal
		}
		if opts.Unmarshal == nil {
			slf._unmarshal = json.Unmarshal
		}
		slf._prefix = opts.Prefix

		pool := &redis.Pool{
			MaxActive:   opts.MaxActive,
			MaxIdle:     opts.MaxIdle,
			IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,

			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial(opts.Network, opts.Addr)
				if err != nil {
					return nil, err
				}
				if opts.Password != "" {
					if _, err := conn.Do("AUTH", opts.Password); err != nil {
						conn.Close()
						return nil, err
					}
				}
				if _, err := conn.Do("SELECT", opts.Db); err != nil {
					conn.Close()
					return nil, err
				}
				return conn, err
			},

			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				_, err := conn.Do("PING")
				return err
			},
		}

		slf._c = pool

		return nil
	default:
		return errors.New("Unsupported options")
	}
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

// Get 获取键值。
func (slf *RedisDB) Get(key string) (interface{}, error) {
	return slf.Do("GET", slf.getKey(key))
}

// GetString 获取string类型的键值
func (slf *RedisDB) GetString(key string) (string, error) {
	return redis.String(slf.Get(key))
}

//GetInt 获取int类型的键值
func (slf *RedisDB) GetInt(key string) (int, error) {
	return redis.Int(slf.Get(key))
}

// GetInt64 获取int64类型的键值
func (slf *RedisDB) GetInt64(key string) (int64, error) {
	return redis.Int64(slf.Get(key))
}

// GetBool 获取bool类型的键值
func (slf *RedisDB) GetBool(key string) (bool, error) {
	return redis.Bool(slf.Get(key))
}

// GetObject 获取非基本类型stuct的键值。在实现上，使用json的Marshal和Unmarshal做序列化存取。
func (slf *RedisDB) GetObject(key string, val interface{}) error {
	reply, err := slf.Get(key)
	return slf.decode(reply, err, val)
}

//Set doc
//@Summary set value
//@param key
//@param value
//@param expire
func (slf *RedisDB) Set(key string, val interface{}, expire int64) error {
	value, err := slf.encode(val)
	if err != nil {
		return err
	}
	if expire > 0 {
		_, err := slf.Do("SETEX", slf.getKey(key), expire, value)
		return err
	}
	_, err = slf.Do("SET", slf.getKey(key), value)
	return err
}

// Exists 检查键是否存在
func (slf *RedisDB) Exists(key string) (bool, error) {
	return redis.Bool(slf.Do("EXISTS", slf.getKey(key)))
}

//Del 删除键
func (slf *RedisDB) Del(key string) error {
	_, err := slf.Do("DEL", slf.getKey(key))
	return err
}

// Flush 清空当前数据库中的所有 key，慎用！
func (slf *RedisDB) Flush() error {
	_, err := slf.Do("FLUSHDB")
	return err
}

// TTL 以秒为单位。当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1
func (slf *RedisDB) TTL(key string) (ttl int64, err error) {
	return redis.Int64(slf.Do("TTL", slf.getKey(key)))
}

// Expire 设置键过期时间，expire的单位为秒
func (slf *RedisDB) Expire(key string, expire int64) error {
	_, err := redis.Bool(slf.Do("EXPIRE", slf.getKey(key), expire))
	return err
}

// Incr 将 key 中储存的数字值增一
func (slf *RedisDB) Incr(key string) (val int64, err error) {
	return redis.Int64(slf.Do("INCR", slf.getKey(key)))
}

// IncrBy 将 key 所储存的值加上给定的增量值（increment）。
func (slf *RedisDB) IncrBy(key string, amount int64) (val int64, err error) {
	return redis.Int64(slf.Do("INCRBY", slf.getKey(key), amount))
}

// Decr 将 key 中储存的数字值减一。
func (slf *RedisDB) Decr(key string) (val int64, err error) {
	return redis.Int64(slf.Do("DECR", slf.getKey(key)))
}

// DecrBy key 所储存的值减去给定的减量值（decrement）。
func (slf *RedisDB) DecrBy(key string, amount int64) (val int64, err error) {
	return redis.Int64(slf.Do("DECRBY", slf.getKey(key), amount))
}

// HMSet 将一个map存到Redis hash，同时设置有效期，单位：秒
// Example:
//
// ```golang
// m := make(map[string]interface{})
// m["name"] = "corel"
// m["age"] = 23
// err := slf.HMSet("user", m, 10)
// ```
func (slf *RedisDB) HMSet(key string, val interface{}, expire int) (err error) {
	conn := slf._c.Get()
	defer conn.Close()
	err = conn.Send("HMSET", redis.Args{}.Add(slf.getKey(key)).AddFlat(val)...)
	if err != nil {
		return
	}
	if expire > 0 {
		err = conn.Send("EXPIRE", slf.getKey(key), int64(expire))
	}
	if err != nil {
		return
	}
	conn.Flush()
	_, err = conn.Receive()
	return
}

/** Redis hash 是一个string类型的field和value的映射表，hash特别适合用于存储对象。 **/

// HSet 将哈希表 key 中的字段 field 的值设为 val
// Example:
//
// ```golang
// _, err := slf.HSet("user", "age", 23)
// ```
func (slf *RedisDB) HSet(key, field string, val interface{}) (interface{}, error) {
	value, err := slf.encode(val)
	if err != nil {
		return nil, err
	}
	return slf.Do("HSET", slf.getKey(key), field, value)
}

// HGet 获取存储在哈希表中指定字段的值
// Example:
//
// ```golang
// val, err := slf.HGet("user", "age")
// ```
func (slf *RedisDB) HGet(key, field string) (reply interface{}, err error) {
	reply, err = slf.Do("HGET", slf.getKey(key), field)
	return
}

// HGetString HGet的工具方法，当字段值为字符串类型时使用
func (slf *RedisDB) HGetString(key, field string) (reply string, err error) {
	reply, err = redis.String(slf.HGet(key, field))
	return
}

// HGetInt HGet的工具方法，当字段值为int类型时使用
func (slf *RedisDB) HGetInt(key, field string) (reply int, err error) {
	reply, err = redis.Int(slf.HGet(key, field))
	return
}

// HGetInt64 HGet的工具方法，当字段值为int64类型时使用
func (slf *RedisDB) HGetInt64(key, field string) (reply int64, err error) {
	reply, err = redis.Int64(slf.HGet(key, field))
	return
}

// HGetBool HGet的工具方法，当字段值为bool类型时使用
func (slf *RedisDB) HGetBool(key, field string) (reply bool, err error) {
	reply, err = redis.Bool(slf.HGet(key, field))
	return
}

// HGetObject HGet的工具方法，当字段值为非基本类型的stuct时使用
func (slf *RedisDB) HGetObject(key, field string, val interface{}) error {
	reply, err := slf.HGet(key, field)
	return slf.decode(reply, err, val)
}

// HGetAll HGetAll("key", &val)
func (slf *RedisDB) HGetAll(key string, val interface{}) error {
	v, err := redis.Values(slf.Do("HGETALL", slf.getKey(key)))
	if err != nil {
		return err
	}

	if err := redis.ScanStruct(v, val); err != nil {
		fmt.Println(err)
	}

	return err
}

/**
Redis列表是简单的字符串列表，按照插入顺序排序。你可以添加一个元素到列表的头部（左边）或者尾部（右边）
**/

// BLPop 它是 LPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 超时参数 timeout 接受一个以秒为单位的数字作为值。超时参数设为 0 表示阻塞时间可以无限期延长(block indefinitely) 。
func (slf *RedisDB) BLPop(key string, timeout int) (interface{}, error) {
	values, err := redis.Values(slf.Do("BLPOP", slf.getKey(key), timeout))
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("redisgo: unexpected number of values, got %d", len(values))
	}
	return values[1], err
}

// BLPopInt BLPop的工具方法，元素类型为int时
func (slf *RedisDB) BLPopInt(key string, timeout int) (int, error) {
	return redis.Int(slf.BLPop(key, timeout))
}

// BLPopInt64 BLPop的工具方法，元素类型为int64时
func (slf *RedisDB) BLPopInt64(key string, timeout int) (int64, error) {
	return redis.Int64(slf.BLPop(key, timeout))
}

// BLPopString BLPop的工具方法，元素类型为string时
func (slf *RedisDB) BLPopString(key string, timeout int) (string, error) {
	return redis.String(slf.BLPop(key, timeout))
}

// BLPopBool BLPop的工具方法，元素类型为bool时
func (slf *RedisDB) BLPopBool(key string, timeout int) (bool, error) {
	return redis.Bool(slf.BLPop(key, timeout))
}

// BLPopObject BLPop的工具方法，元素类型为object时
func (slf *RedisDB) BLPopObject(key string, timeout int, val interface{}) error {
	reply, err := slf.BLPop(key, timeout)
	return slf.decode(reply, err, val)
}

// BRPop 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 超时参数 timeout 接受一个以秒为单位的数字作为值。超时参数设为 0 表示阻塞时间可以无限期延长(block indefinitely) 。
func (slf *RedisDB) BRPop(key string, timeout int) (interface{}, error) {
	values, err := redis.Values(slf.Do("BRPOP", slf.getKey(key), timeout))
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("redisgo: unexpected number of values, got %d", len(values))
	}
	return values[1], err
}

// BRPopInt BRPop的工具方法，元素类型为int时
func (slf *RedisDB) BRPopInt(key string, timeout int) (int, error) {
	return redis.Int(slf.BRPop(key, timeout))
}

// BRPopInt64 BRPop的工具方法，元素类型为int64时
func (slf *RedisDB) BRPopInt64(key string, timeout int) (int64, error) {
	return redis.Int64(slf.BRPop(key, timeout))
}

// BRPopString BRPop的工具方法，元素类型为string时
func (slf *RedisDB) BRPopString(key string, timeout int) (string, error) {
	return redis.String(slf.BRPop(key, timeout))
}

// BRPopBool BRPop的工具方法，元素类型为bool时
func (slf *RedisDB) BRPopBool(key string, timeout int) (bool, error) {
	return redis.Bool(slf.BRPop(key, timeout))
}

// BRPopObject BRPop的工具方法，元素类型为object时
func (slf *RedisDB) BRPopObject(key string, timeout int, val interface{}) error {
	reply, err := slf.BRPop(key, timeout)
	return slf.decode(reply, err, val)
}

// LPop 移出并获取列表中的第一个元素（表头，左边）
func (slf *RedisDB) LPop(key string) (interface{}, error) {
	return slf.Do("LPOP", slf.getKey(key))
}

// LPopInt 移出并获取列表中的第一个元素（表头，左边），元素类型为int
func (slf *RedisDB) LPopInt(key string) (int, error) {
	return redis.Int(slf.LPop(key))
}

// LPopInt64 移出并获取列表中的第一个元素（表头，左边），元素类型为int64
func (slf *RedisDB) LPopInt64(key string) (int64, error) {
	return redis.Int64(slf.LPop(key))
}

// LPopString 移出并获取列表中的第一个元素（表头，左边），元素类型为string
func (slf *RedisDB) LPopString(key string) (string, error) {
	return redis.String(slf.LPop(key))
}

// LPopBool 移出并获取列表中的第一个元素（表头，左边），元素类型为bool
func (slf *RedisDB) LPopBool(key string) (bool, error) {
	return redis.Bool(slf.LPop(key))
}

// LPopObject 移出并获取列表中的第一个元素（表头，左边），元素类型为非基本类型的struct
func (slf *RedisDB) LPopObject(key string, val interface{}) error {
	reply, err := slf.LPop(key)
	return slf.decode(reply, err, val)
}

// RPop 移出并获取列表中的最后一个元素（表尾，右边）
func (slf *RedisDB) RPop(key string) (interface{}, error) {
	return slf.Do("RPOP", slf.getKey(key))
}

// RPopInt 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int
func (slf *RedisDB) RPopInt(key string) (int, error) {
	return redis.Int(slf.RPop(key))
}

// RPopInt64 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int64
func (slf *RedisDB) RPopInt64(key string) (int64, error) {
	return redis.Int64(slf.RPop(key))
}

// RPopString 移出并获取列表中的最后一个元素（表尾，右边），元素类型为string
func (slf *RedisDB) RPopString(key string) (string, error) {
	return redis.String(slf.RPop(key))
}

// RPopBool 移出并获取列表中的最后一个元素（表尾，右边），元素类型为bool
func (slf *RedisDB) RPopBool(key string) (bool, error) {
	return redis.Bool(slf.RPop(key))
}

// RPopObject 移出并获取列表中的最后一个元素（表尾，右边），元素类型为非基本类型的struct
func (slf *RedisDB) RPopObject(key string, val interface{}) error {
	reply, err := slf.RPop(key)
	return slf.decode(reply, err, val)
}

// LPush 将一个值插入到列表头部
func (slf *RedisDB) LPush(key string, member interface{}) error {
	value, err := slf.encode(member)
	if err != nil {
		return err
	}
	_, err = slf.Do("LPUSH", slf.getKey(key), value)
	return err
}

// RPush 将一个值插入到列表尾部
func (slf *RedisDB) RPush(key string, member interface{}) error {
	value, err := slf.encode(member)
	if err != nil {
		return err
	}
	_, err = slf.Do("RPUSH", slf.getKey(key), value)
	return err
}

// LREM 根据参数 count 的值，移除列表中与参数 member 相等的元素。
// count 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 member 相等的元素，数量为 count 。
// count < 0 : 从表尾开始向表头搜索，移除与 member 相等的元素，数量为 count 的绝对值。
// count = 0 : 移除表中所有与 member 相等的值。
// 返回值：被移除元素的数量。
func (slf *RedisDB) LREM(key string, count int, member interface{}) (int, error) {
	return redis.Int(slf.Do("LREM", slf.getKey(key), count, member))
}

// LLen 获取列表的长度
func (slf *RedisDB) LLen(key string) (int64, error) {
	return redis.Int64(slf.Do("RPOP", slf.getKey(key)))
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 和编程语言区间函数的区别：end 下标也在 LRANGE 命令的取值范围之内(闭区间)。
func (slf *RedisDB) LRange(key string, start, end int) (interface{}, error) {
	return slf.Do("LRANGE", slf.getKey(key), start, end)
}

/**
Redis 有序集合和集合一样也是string类型元素的集合,且不允许重复的成员。
不同的是每个元素都会关联一个double类型的分数。redis正是通过分数来为集合中的成员进行从小到大的排序。
有序集合的成员是唯一的,但分数(score)却可以重复。
集合是通过哈希表实现的，所以添加，删除，查找的复杂度都是O(1)。
**/

// ZAdd 将一个 member 元素及其 score 值加入到有序集 key 当中。
func (slf *RedisDB) ZAdd(key string, score int64, member string) (reply interface{}, err error) {
	return slf.Do("ZADD", slf.getKey(key), score, member)
}

// ZRem 移除有序集 key 中的一个成员，不存在的成员将被忽略。
func (slf *RedisDB) ZRem(key string, member string) (reply interface{}, err error) {
	return slf.Do("ZREM", slf.getKey(key), member)
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (slf *RedisDB) ZScore(key string, member string) (int64, error) {
	return redis.Int64(slf.Do("ZSCORE", slf.getKey(key), member))
}

// ZRank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (slf *RedisDB) ZRank(key, member string) (int64, error) {
	return redis.Int64(slf.Do("ZRANK", slf.getKey(key), member))
}

// ZRevrank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (slf *RedisDB) ZRevrank(key, member string) (int64, error) {
	return redis.Int64(slf.Do("ZREVRANK", slf.getKey(key), member))
}

// ZRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (slf *RedisDB) ZRange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(slf.Do("ZRANGE", slf.getKey(key), from, to, "WITHSCORES"))
}

// ZRevrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (slf *RedisDB) ZRevrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(slf.Do("ZREVRANGE", slf.getKey(key), from, to, "WITHSCORES"))
}

// ZRangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (slf *RedisDB) ZRangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(slf.Do("ZRANGEBYSCORE", slf.getKey(key), from, to, "WITHSCORES", "LIMIT", offset, count))
}

// ZRevrangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (slf *RedisDB) ZRevrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(slf.Do("ZREVRANGEBYSCORE", slf.getKey(key), from, to, "WITHSCORES", "LIMIT", offset, count))
}

//Close doc
//@Summary close redis db operation
//@Method Close
func (slf *RedisDB) Close() {
	slf._c.Close()
}

//getKey 将健名加上指定的前缀。
func (slf *RedisDB) getKey(key string) string {
	return slf._prefix + key
}

// encode 序列化要保存的值
func (slf *RedisDB) encode(val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := slf._marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return value, nil
}

// decode 反序列化保存的struct对象
func (slf *RedisDB) decode(reply interface{}, err error, val interface{}) error {
	str, err := redis.String(reply, err)
	if err != nil {
		return err
	}
	return slf._unmarshal([]byte(str), val)
}
