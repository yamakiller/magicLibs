package contractor

import (
	"crypto/rc4"
	"encoding/binary"
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/yamakiller/magicLibs/encryption/dh64"
	"github.com/yamakiller/magicLibs/util"
)

type state int32

const (
	stateIdle     = state(0)
	stateReserved = state(1)
	stateConfirm  = state(2)
)

//Snk 序列ID生成器
type Snk struct {
	dh64.KeyExchange
	HandShakeTime int64
	_sns          snks
}

//Reserve 预约ID
func (slf *Snk) Reserve(pubExKey uint64) (uint32, uint64, error) {
	prvKey, pubKey := slf.KeyPair()
	secret := slf.Secret(prvKey, pubExKey)
	slf._sns.Lock()
	s, err := slf._sns.Next()
	if err != nil {
		slf._sns.Unlock()
		return 0, 0, err
	}
	s._death = util.ToTimestamp(time.Now().
		Add(time.Duration(slf.HandShakeTime) * time.Millisecond))
	s._secret = secret
	tmpBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(tmpBuf, s._secret)

	sc, err := rc4.NewCipher(tmpBuf)
	if err != nil {
		slf._sns.Remove(s._conv)
		slf._sns.Unlock()
		return 0, 0, err
	}
	s._rc4 = sc
	s._state = stateReserved
	slf._sns.Unlock()

	binary.BigEndian.PutUint32(tmpBuf, s._conv)
	s._rc4.XORKeyStream(tmpBuf[4:], tmpBuf[:4])

	return binary.BigEndian.Uint32(tmpBuf[4:]), pubKey, nil
}

//Confirm 确认
func (slf *Snk) Confirm(d []byte) bool {
	if len(d) != 8 {
		return false
	}

	id := binary.BigEndian.Uint32(d)

	slf._sns.Lock()
	s := slf._sns.Get(id)
	if s == nil {
		slf._sns.Unlock()
		return false
	}

	slf._sns.Unlock()

	prc4 := uintptr(unsafe.Pointer(s._rc4))
	p := atomic.LoadPointer(&prc4)

	if p == nil {
		return false
	}

	if s._rc4 == nil || s._state == stateIdle {
		slf._sns.Unlock()
		return false
	}

	tmpBuf := make([]byte, 4)
	s._rc4.XORKeyStream(tmpBuf, d[4:])
	token := binary.BigEndian.Uint32(tmpBuf)
	if token != id {
		slf._sns.Unlock()
		return false
	}

	s._state = stateConfirm

	return true
}

//Authorized 是否授权
func (slf *Snk) Authorized(id uint32) bool {
	slf._sns.Lock()
	s := slf._sns.Get(id)
	if s == nil {
		slf._sns.Unlock()
		return false
	}

	if s._state != stateConfirm {
		slf._sns.Unlock()
		return false
	}
	slf._sns.Unlock()
	return true
}

//Close 关闭
func (slf *Snk) Close(id uint32) {
	slf._sns.Lock()
	slf._sns.Remove(id)
	slf._sns.Unlock()
}

//Update 维护
func (slf *Snk) Update() {

	for {

	}
}

type snks struct {
	_mask  uint32
	_max   uint32 //2的幂
	_snint uint32
	_sz    int
	_cap   uint32
	_ss    []snkData
	sync.Mutex
}

func (slf *snks) Next() (*snkData, error) {
	var i uint32
	for {

		for i = 0; i < slf._cap; i++ {
			key := ((i + slf._snint) & slf._mask)
			if key == 0 {
				key = 1
			}
			hash := key & (slf._cap - 1)
			if slf._ss[hash]._conv == 0 {
				slf._snint = key + 1
				slf._ss[hash]._conv = key
				slf._sz++
				return &slf._ss[hash], nil
			}
		}

		newCap := slf._cap * 2
		if newCap > slf._max {
			newCap = slf._max
		}

		if newCap == slf._cap {
			return nil, errors.New("full")
		}

		slf._ss = append(slf._ss, make([]snkData, newCap-slf._cap)...)
		for i = 0; i < slf._cap; i++ {
			if slf._ss[i]._conv == 0 {
				continue
			}

			hash := slf._ss[i]._conv & uint32(newCap-1)
			if hash == i {
				continue
			}

			tmp := slf._ss[i]
			slf._ss[hash] = tmp
			slf._ss[i]._conv = 0
			slf._ss[i]._death = 0
			slf._ss[i]._secret = 0
			slf._ss[i]._state = stateIdle
		}
		slf._cap = newCap
	}
}

//Get 获取一个元素
func (slf *snks) Get(key uint32) *snkData {
	hash := key & uint32(slf._cap-1)
	if slf._ss[hash]._conv != 0 && slf._ss[hash]._conv == key {
		return &slf._ss[hash]
	}
	return nil
}

//Remove 删除一个元素
func (slf *snks) Remove(key uint32) bool {
	hash := uint32(key) & uint32(slf._cap-1)
	if slf._ss[hash]._conv != 0 && slf._ss[hash]._conv == key {
		slf._ss[hash]._conv = 0
		slf._ss[hash]._death = 0
		slf._ss[hash]._secret = 0
		slf._ss[hash]._state = stateIdle
		slf._sz--
		return true
	}

	return false
}

type snkData struct {
	_conv   uint32
	_secret uint64
	_rc4    *rc4.Cipher
	_death  int64
	_state  state
}
