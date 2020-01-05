package vm

type RaspState int

const (
	DataTag RaspState = 1 + iota
	CalcTag
	AddrTag // address of the executing contract
	OriginAddrTag
	CallerAddrTag
	BalanceTag
)

type RaspTag map[int]RaspState

type RaspTag64 map[uint64]RaspState

func (ST RaspTag) push(loc int, state RaspState) {
	ST[loc] = state
}

func (ST RaspTag) delete(loc int) {
	delete(ST, loc)
}

func (ST RaspTag) check(loc int) bool {
	_, ok := ST[loc]
	return ok
}

// 如果是memory,根据是mstore8h和mstore32决定N的值；如果是storage，N=1
func (ST RaspTag64) push64(loc uint64, N uint64, state RaspState) {
	for i := uint64(0); i < N; i++ {
		ST[i+loc] = state
	}
}

func (ST RaspTag64) delete64(loc uint64, N uint64) {
	for i := uint64(0); i < N; i++ {
		delete(ST, loc)
	}
}

func (ST RaspTag64) check64(loc uint64) bool {
	_, ok := ST[loc]
	return ok
}
