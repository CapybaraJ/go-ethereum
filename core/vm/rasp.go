package vm

import "fmt"

var hookstate = 0
var StackTags RaspTag
var MemTags RaspTag64
var StorageTags RaspTag64

func init() {
	StackTags = make(RaspTag)
	MemTags = make(RaspTag64)
	StorageTags = make(RaspTag64)
}

func InRASP(pc uint64, op OpCode, contract *Contract, input []byte, st *Stack, mem *Memory) bool {
	//fmt.Printf("REVERT!!!!: %v(%v) || the stack is: %X\n", op, pc, stack.data)
	result := true
	switch op {
	case ADD, SUB, MUL, DIV, AND:
		if StackTags.check(st.len()) || StackTags.check(st.len()-1) {
			fmt.Println("effective Caculate")
			result = HookCalc(pc, op, contract, input, st)
			StackTags.delete(st.len())
			StackTags.push(st.len()-1, CalcTag)
		}
		break
	case SDIV, MOD, SMOD, ADDMOD, MULMOD, EXP, SIGNEXTEND:
		break
	default:
		HookVar(op, contract, input, st)
	}
	return result
}

//
func HookVar(op OpCode, contract *Contract, input []byte, st *Stack) {
	switch op {
	case ISZERO, BALANCE, EXTCODESIZE, EXTCODEHASH, POP, JUMP:
		if StackTags.check(st.len()) {
			StackTags.delete(st.len())
		}
		break
	// 0x10 range - comparison ops.
	case LT, GT, SLT, SGT, EQ, SHA3, JUMPI:
		if StackTags.check(st.len()) { //&& StackTags[st.len()] == DataTag
			StackTags.delete(st.len())
		}
		if StackTags.check(st.len() - 1) { //&& StackTags[st.len()-1] == DataTag
			StackTags.delete(st.len() - 1)
		}
		break
	case AND, OR, XOR, BYTE, SHR, SHL, SAR:
		if StackTags.check(st.len()) || StackTags.check(st.len()-1) {
			StackTags.delete(st.len() - 1)
		}
	// 0x30 range - closure state. ignore RETURNDATASIZE CODECOPY CALLDATASIZE CODESIZE CALLVALUE
	case ADDRESS:
		StackTags.push(st.len()+1, AddrTag)
	case ORIGIN:
		StackTags.push(st.len()+1, OriginAddrTag)
	case CALLER:
		StackTags.push(st.len()+1, CallerAddrTag)
	case CALLDATALOAD:
		StackTags.push(st.len(), DataTag)
		break
	case CALLDATACOPY, CODECOPY, EXTCODECOPY: //注意内存操作
		for i := 0; i < 3; i++ {
			StackTags.delete(st.len() - i)
		}
		if op == CALLDATACOPY {
			MemTags.push64(st.peek().Uint64(), st.Back(2).Uint64(), DataTag)
		}
		if op == EXTCODECOPY {
			StackTags.delete(st.len() - 3)
		}
		break
	case MLOAD:
		offset := st.peek().Uint64()
		for i := uint64(0); i < 32; i++ {
			if MemTags.check64(i + offset) {
				StackTags.push(st.len(), MemTags[i+offset])
			}
		}
		break
	case MSTORE, MSTORE8:
		if StackTags.check(st.len() - 1) {
			offset := st.peek().Uint64()
			len := 0
			if op == MSTORE {
				len = 32
			} else {
				len = 8
			}
			MemTags.push64(offset, uint64(len), StackTags[st.len()-1])
		}
		StackTags.delete(st.len())
		StackTags.delete(st.len() - 1)
		break
	case SSTORE:
		if StackTags.check(st.len() - 1) {
			StorageTags.push64(st.peek().Uint64(), uint64(1), StackTags[st.len()-1])
		}
		StackTags.delete(st.len())
		StackTags.delete(st.len() - 1)
		break
	case SLOAD:
		key := st.peek().Uint64()
		if StorageTags.check64(key) {
			StackTags.push(st.len(), StorageTags[key])
		} else {
			StackTags.delete(st.len())
		}
		break
	case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8, DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
		loc := int(op - DUP1)
		ok := StackTags.check(st.len() - loc)
		if ok {
			StackTags.push(st.len()+1, StackTags[st.len()-loc])
			return
		}
		break
	case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8, SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
		loc := int(op - SWAP1)
		ok1 := StackTags.check(st.len())
		ok2 := StackTags.check(st.len() - loc - 1)
		if ok1 && ok2 {
			tmp := StackTags[st.len()]
			StackTags[st.len()] = StackTags[st.len()-loc-1]
			StackTags[st.len()-loc-1] = tmp
			return
		} else if ok1 && !ok2 {
			tmp := StackTags[st.len()]
			delete(StackTags, st.len())
			StackTags.push(st.len()-loc-1, tmp)
		} else if !ok1 && ok2 {
			tmp := StackTags[st.len()-loc-1]
			delete(StackTags, st.len()-loc-1)
			StackTags.push(st.len(), tmp)
		}
		break
	case LOG0, LOG1, LOG2, LOG3, LOG4:
		sum := int(op-LOG0) + 1
		for i := 0; i <= sum; i++ {
			StackTags.delete(st.len() - i)
		}
		break
	// 0xf0 range - closures.
	case CREATE:
		for i := 0; i <= 2; i++ {
			StackTags.delete(st.len() - i)
		}
	case CALL, CALLCODE:
		for i := 0; i <= 6; i++ {
			StackTags.delete(st.len() - i)
		}
	case RETURN, REVERT:
		for i := 0; i <= 1; i++ {
			StackTags.delete(st.len() - i)
		}
	case DELEGATECALL, STATICCALL:
		for i := 0; i <= 5; i++ {
			StackTags.delete(st.len() - i)
		}
	case CREATE2:
		for i := 0; i <= 3; i++ {
			StackTags.delete(st.len() - i)
		}
	case SELFDESTRUCT:
		StackTags.delete(st.len())
	default:
		break
	}
}
