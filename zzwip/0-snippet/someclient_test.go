package snippet

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_SomeClient(t *testing.T) {
	url := "some.com"
	c := NewSome(url)
	if err := c.DoSomething(); err != nil {
		t.Error(err)
	}
}

func Test_StructSize(_ *testing.T) {
	type T struct {
		t1 byte
		t2 int32
		t3 int64
		t4 string
		t5 bool
	}

	t := &T{
		t1: 1,
		t2: 2,
		t3: 3,
		t4: "", // In golang, the size of string is 16 bytes for any strings (even empty).
		t5: true,
	}

	// Strings in Go are represented by reflect.StringHeader containing
	// a pointer (uintptr 8bytes) to actual string data
	// and a length (int 8bytes) of string:

	fmt.Printf("t1 offset: %d, size: %d\n", unsafe.Offsetof(t.t1), unsafe.Sizeof(t.t1))
	fmt.Printf("t2 offset: %d, size: %d\n", unsafe.Offsetof(t.t2), unsafe.Sizeof(t.t2))
	fmt.Printf("t3 offset: %d, size: %d\n", unsafe.Offsetof(t.t3), unsafe.Sizeof(t.t3))
	fmt.Printf("t4 offset: %d, size: %d\n", unsafe.Offsetof(t.t4), unsafe.Sizeof(t.t4))
	fmt.Printf("t5 offset: %d, size: %d\n", unsafe.Offsetof(t.t5), unsafe.Sizeof(t.t5))

	fmt.Printf("t struct size: %d\n", unsafe.Sizeof(*t))
	fmt.Printf("t point size: %d\n", unsafe.Sizeof(t))

	a := 0b01000000000010000110100100010000

	fmt.Println(a)
	fmt.Printf("%d\n", a)
	fmt.Printf("%#x\n", a)

	fmt.Printf("int: %d\n", unsafe.Sizeof(int(0)))
	fmt.Printf("int32: %d\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int64: %d\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("slice empty: %d\n", unsafe.Sizeof([]uint8{}))
	fmt.Printf("slice 5: %d\n", unsafe.Sizeof([]uint8{1, 2, 3, 4, 5}))
	fmt.Printf("slice 6: %d\n", unsafe.Sizeof([]uint8{1, 2, 3, 4, 5, 6}))
	fmt.Printf("unsafe.Pointer: %d\n", unsafe.Sizeof(unsafe.Pointer(&a)))

}
