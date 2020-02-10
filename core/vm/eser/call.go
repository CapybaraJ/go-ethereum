package eser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

var CallChain = NewLinkedList()

func ReentryCheck() ErrorType {
	if hasCycle(CallChain) {
		return ReentryError
	}
	return NoError
}

/**
 * 判断链表是否成环
 */
func hasCycle(list *LinkedList) bool {
	head := list.Head
	if head == nil {
		return false
	}
	t, h := head, head
	started := false
	for h != nil && h.next != nil {
		value1, ok1 := h.X.([]byte)
		value2, ok2 := t.X.([]byte)
		fmt.Println(value1, ok1)
		fmt.Println(value2, ok2)
		result := false
		if ok1 && ok2 {
			result = bytes.Equal(value1, value2)
		} else {
			if t == h {
				result = true
			}
		}
		if result {
			if started {
				return true
			} else {
				started = true
			}
		}
		t = t.next
		h = h.next.next
	}
	return false
}

func GetSHA1(str []byte) []byte {
	//SHA1
	Sha1Inst := sha1.New()
	Sha1Inst.Write(str)
	result := Sha1Inst.Sum([]byte(""))
	return result
}
