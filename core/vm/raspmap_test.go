package vm

import (
	"fmt"
	"testing"
)

func TestIt(t *testing.T) {
	fmt.Println("Starting test...")
	ml := make(RaspTag)
	ml.push(99, CalcTag)
	ml.push(7, DataTag)
	ml.push(17, DataTag)
	ml.push(37, DataTag)
	fmt.Printf("%v\n", ml)
	ml.delete(17)
	fmt.Print(ml.check(66))
	fmt.Printf("%v\n", ml)
}
