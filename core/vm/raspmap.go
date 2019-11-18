package vm

type RaspState byte

const (
	DataTag RaspState = iota
	CalcTag
)

type RaspTag map[int]RaspState

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
