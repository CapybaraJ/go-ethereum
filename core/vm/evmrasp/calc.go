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
		break
	case -1:
		return allFF.check(val, mid+1, right)
		break
	case 0:
		return mid
		break
	}
	return -1
}

//
func HookAdd(a *big.Int, b *big.Int, len int) bool {
	c := big.NewInt(int64(0))
	c.Add(a, b)
	fmt.Println("CMP1: ", c, MyMaxBig256)
	if c.Cmp(MyMaxBig256) == 1 {
		//上溢出
		fmt.Println("上溢出")
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
