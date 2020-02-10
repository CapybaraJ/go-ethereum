package vm

import (
	"github.com/ethereum/go-ethereum/core/vm/eser"
)

func HookCalc(pc uint64, op OpCode, contract *Contract, input []byte, st *Stack) eser.ErrorType {
	result := eser.NoError
	switch op {
	case ADD:
		result = eser.HookAdd(st.peek(), st.Back(1), st.len())
		break
	case SUB:
		result = eser.HookSub(st.peek(), st.Back(1), st.len())
		break
	case MUL:
		result = eser.HookMul(st.peek(), st.Back(1), st.len())
		break
	case AND:
		if StackTags.check(st.len()-1) && StackTags[st.len()-1] == CalcTag {
			if eser.HookAnd(st.peek(), st.Back(1)) == eser.NoError {
				StackTags[st.len()-1] = DataTag
			}
		}
	default:
		break
	}
	return result
}
