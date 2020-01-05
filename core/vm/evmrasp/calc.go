package evmrasp

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

type FFs [31]*big.Int

var allFF FFs

var (
	my256       = math.BigPow(2, 256)
	my256m1     = new(big.Int).Sub(my256, big.NewInt(1))
	MyMaxBig256 = new(big.Int).Set(my256m1)
)

func init() {
	var ffs bytes.Buffer
	for a := 0; a < 31; a++ {
		ffs.WriteString("ff")
		allFF[a], _ = new(big.Int).SetString(ffs.String(), 16)
	}
	//fmt.Printf("len:%v, %v", len(allFF), allFF)
}

// 计算出到底是AND 0xfffff几个f决定是 uint8 uint16 ... uint 248
func (FFs) check(val *big.Int, left int, right int) int {
	if left == right {
		return -1
	}
	mid := (left + right) / 2
	switch allFF[mid].Cmp(val) {
	case 1:
		return allFF.check(val, left, mid)
	case -1:
		return allFF.check(val, mid+1, right)
	case 0:
		return mid
	}
	return -1
}

//
func HookAdd(a *big.Int, b *big.Int, len int) bool {
	c := big.NewInt(int64(0))
	c.Add(a, b)
	fmt.Println("CMP1: ", c, MyMaxBig256)
	if c.Cmp(MyMaxBig256) == 1 {
		fmt.Println("上溢出") //上溢出
		return false
	}
	return true
}

func HookAnd(a *big.Int, b *big.Int) bool {
	loc := allFF.check(a, 0, 31)
	fmt.Print("loc:", loc)
	if loc != -1 {
		// get the result
		fmt.Println("CMP2: ", b, allFF[loc])
		if b.Cmp(allFF[loc]) == 1 {
			// 超过限制，上溢出
			fmt.Println("AND 上溢出", b, a)
			return false
		}
	}
	return true
}

func HookSub(a *big.Int, b *big.Int, len int) bool {
	fmt.Println("CMP2: ", a, b)
	if b.Cmp(a) == 1 {
		fmt.Println("下溢出") //上溢出
		return false
	}
	return true
}

func HookMul(a *big.Int, b *big.Int, len int) bool {
	c := big.NewInt(int64(0))
	d := big.NewInt(int64(1))
	if a.Cmp(c) == 0 || b.Cmp(c) == 0 || a.Cmp(d) == 0 || b.Cmp(d) == 0 {
		return true
	}
	c.Mul(a, b)
	fmt.Println("CMP3: ", c)
	if c.Cmp(MyMaxBig256) == 1 {
		fmt.Println("上溢出") //上溢出
		return false
	}
	return true
}
