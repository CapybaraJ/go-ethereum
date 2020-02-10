package eser

type ErrorType int

const (
	NoError ErrorType = 0xeeee + iota
	CalcError
	ReentryError
	DOSError // address of the executing contract
	DelegateError
)
