package dh64

import (
	"math/rand"
)

const (
	p uint64 = 0xffffffffffffffc5
	g uint64 = 5
)

func mul_mod_p(a, b uint64) uint64 {
	var m uint64 = 0
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

func pow_mod_p(a, b uint64) uint64 {
	if b == 1 {
		return a
	}
	t := pow_mod_p(a, b>>1)
	t = mul_mod_p(t, t)
	if b%2 > 0 {
		t = mul_mod_p(t, a)
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
	return pow_mod_p(a, b)
}

func KeyPair() (privateKey, publicKey uint64) {
	a := uint64(rand.Uint32())
	b := uint64(rand.Uint32()) + 1
	privateKey = (a << 32) | b
	publicKey = PublicKey(privateKey)
	return
}

func PublicKey(privateKey uint64) uint64 {
	return powmodp(g, privateKey)
}

func Secret(privateKey, anotherPublicKey uint64) uint64 {
	return powmodp(anotherPublicKey, privateKey)
}

//
//1.随机生成一对64位密钥（私钥 + 公钥) myPrivateKey, myPublicKey := dh64.KeyPair()
//2.公锁发送给客户端
//3.等待客户端的公锁
//4.根据客户端的公锁 + 服务器的私锁，计算出密钥：secert := dh64.Secert(myPrivateKey, anotherPublicKey);
//5.客户端既按照此方法计算出密钥
