package vm

import (
	"github.com/ethereum/go-ethereum/core/vm/evmrasp"
)

func HookCalc(pc uint64, op OpCode, contract *Contract, input []byte, st *Stack) {
	switch op {
	case ADD:
		if evmrasp.HookAdd(st.peek(), st.Back(1), st.len()) {
			StackTags.push(st.len()-1, CalcTag)
		}
		break
	case SUB:
	case AND:
		if StackTags.check(st.len()-1) && StackTags[st.len()-1] == CalcTag {
			if evmrasp.HookAnd(st.peek(), st.Back(1)) {
				StackTags[st.len()-1] = DataTag
			}
		}
	}
}
