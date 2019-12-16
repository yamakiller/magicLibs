package dh64

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
)

const (
	p uint64 = 0xffffffffffffffc5
	g uint64 = 5
)

func mulModP(a, b uint64) uint64 {
	var m uint64
	for b > 0 {
		if b&1 > 0 {
			t := p - a
			if m >= t {
				m -= t
			} else {
				m += a
			}
		}
		if a >= p-a {
			a = a*2 - p
		} else {
			a = a * 2
		}
		b >>= 1
	}
	return m
}

func powModP(a, b uint64) uint64 {
	if b == 1 {
		return a
	}
	t := powModP(a, b>>1)
	t = mulModP(t, t)
	if b%2 > 0 {
		t = mulModP(t, a)
	}
	return t
}

func powmodp(a uint64, b uint64) uint64 {
	if a == 0 {
		panic("DH64 zero public key")
	}
	if b == 0 {
		panic("DH64 zero private key")
	}
	if a > p {
		a %= p
	}
	return powModP(a, b)
}

//KeyPair Generate public key key pair
func KeyPair() (privateKey, publicKey uint64) {
	a := uint64(rand.Uint32())
	b := uint64(rand.Uint32()) + 1
	privateKey = (a << 32) | b
	publicKey = PublicKey(privateKey)
	return
}

//PublicKey private key to public key
func PublicKey(privateKey uint64) uint64 {
	return powmodp(g, privateKey)
}

//Secret Generate Secret to uint64
func Secret(privateKey, anotherPublicKey uint64) uint64 {
	return powmodp(anotherPublicKey, privateKey)
}

//SecretToString Generate Secret to String
func SecretToString(privateKey, anotherPublicKey uint64) string {
	secret := Secret(privateKey, anotherPublicKey)
	tmpBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tmpBytes, secret)
	return hex.EncodeToString(tmpBytes)
}

//
//1.随机生成一对64位密钥（私钥 + 公钥) myPrivateKey, myPublicKey := dh64.KeyPair()
//2.公锁发送给客户端
//3.等待客户端的公锁
//4.根据客户端的公锁 + 服务器的私锁，计算出密钥：secert := dh64.Secert(myPrivateKey, anotherPublicKey);
//5.客户端既按照此方法计算出密钥
