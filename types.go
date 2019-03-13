package main

import "fmt"

type OutOfBondsError struct {
	errOffset int
	lenSlice  int
}

func (e OutOfBondsError) Error() string {
	return fmt.Sprintf("Offset is %d and slice has only %d elements", e.errOffset, e.lenSlice)
}
