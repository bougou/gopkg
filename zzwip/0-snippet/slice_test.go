package snippet

import (
	"fmt"
	"testing"
)

// see: https://medium.com/swlh/golang-tips-why-pointers-to-slices-are-useful-and-how-ignoring-them-can-lead-to-tricky-bugs-cac90f72e77b
// Therefore make sure to keep in mind that you can pass a slice by value
// if you want to modify merely the values of the elements, not their number or position,
// otherwise weird bugs will arise from time to time.

func Test_Slice1(t *testing.T) {
	printSlice1 := func(slice []string) {
		slice[0] = "b"
		slice[1] = "b"
		fmt.Println(slice)
	}

	slice := []string{"a", "a"}
	printSlice1(slice) // [b b]
	fmt.Println(slice) // [b b]
}

func Test_Slice2(t *testing.T) {
	printSlice2 := func(slice *[]string) {
		(*slice)[0] = "b"
		(*slice)[1] = "b"
		fmt.Println(*slice)
	}

	slice := []string{"a", "a"}
	printSlice2(&slice) // [b b]
	fmt.Println(slice)  // [b b]
}

func Test_Slice3(t *testing.T) {
	printSlice1 := func(slice []string) {
		// Warning,
		// Here the slice gets allocated, a new location of the memory is used.
		slice = append(slice, "c")

		// modification happens on new memory location.
		slice[0] = "b"
		slice[1] = "b"

		fmt.Println(slice)
	}

	slice := []string{"a", "a"}
	printSlice1(slice) // [b b c]
	fmt.Println(slice) // [a a]
}

func Test_Slice4(t *testing.T) {
	printSlice1 := func(slice *[]string) {
		// Warning,
		// Here the slice gets allocated, a new location of the memory is used.
		// By using pointer, the origin slice now points to this new memory.
		*slice = append(*slice, "c")

		// modification happens on new memory location.
		(*slice)[0] = "b"
		(*slice)[1] = "b"

		fmt.Println(*slice)
	}

	slice := []string{"a", "a"}
	printSlice1(&slice) // [b b c]
	fmt.Println(slice)  // [b b c]
}
