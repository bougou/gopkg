package snippet

import (
	"fmt"
	"testing"
)

type SomeStruct struct{}

// 1. Method receivers default to pointers except in rare cases.
func (b *SomeStruct) Method() {}

// 2. Slices, maps, channels, strings, function values, and interface values are
// implemented with pointers internally, and a pointer to them is often redundant.

func Test_Receiver(t *testing.T) {
	a := []*SomeStruct{}

	b := SomeStruct{}
	a = append(a, &b)

	fmt.Println(a)

}
