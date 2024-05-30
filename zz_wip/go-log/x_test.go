package main

import (
	"fmt"
	"reflect"
	"testing"
)

// https://stackoverflow.com/questions/59040161/given-a-method-value-get-receiver-object

type MyStruct struct {
	Inner string
}

func (m MyStruct) MyMethod(a string) error {
	fmt.Printf("my method, inner: (%s), params: (%s)\n", m.Inner, a)
	return nil
}

func SpecialFunc(fn func(a string) error) {
	fmt.Println(reflect.TypeOf(fn))
	_ = fn("1")
}

func Test_x(t *testing.T) {
	a := MyStruct{
		Inner: "hello",
	}
	_ = a.MyMethod("a")

	SpecialFunc(a.MyMethod)
}
