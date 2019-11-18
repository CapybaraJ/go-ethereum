package vm

//var calcrasp = CalcResult{
//Result:    big.NewInt(int64(0)),
//State:	0,
//Overflow:  false,
//Underflow: false,
//}

var hookstate = 0

func InRASP(pc uint64, op OpCode, contract *Contract, input []byte, st *Stack, mem *Memory) {
	//fmt.Printf("REVERT!!!!: %v(%v) || the stack is: %X\n", op, pc, stack.data)
	switch op {
	case ADD, AND:
		if StackTags.check(st.len()) || StackTags.check(st.len()-1) {
			HookCalc(pc, op, contract, input, st)
		}
		break
	default:
		HookVar(op, contract, input, st)
	}
}

//
func HookVar(op OpCode, contract *Contract, input []byte, st *Stack) {
	switch op {
	case CALLDATALOAD:
		StackTags.push(st.len(), DataTag)
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
	}
}
