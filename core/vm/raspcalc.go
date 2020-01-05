package vm

import (
	"github.com/ethereum/go-ethereum/core/vm/evmrasp"
)

func HookCalc(pc uint64, op OpCode, contract *Contract, input []byte, st *Stack) bool {
	result := true
	switch op {
	case ADD:
		result = evmrasp.HookAdd(st.peek(), st.Back(1), st.len())
		break
	case SUB:
		result = evmrasp.HookSub(st.peek(), st.Back(1), st.len())
		break
	case MUL:
		result = evmrasp.HookMul(st.peek(), st.Back(1), st.len())
		break
	case AND:
		if StackTags.check(st.len()-1) && StackTags[st.len()-1] == CalcTag {
			if evmrasp.HookAnd(st.peek(), st.Back(1)) {
				StackTags[st.len()-1] = DataTag
			}
		}
	default:
		break
	}
	return result
}
