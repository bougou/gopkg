package snippet

import (
	"bytes"
	"fmt"
	"testing"
)

// Also it's interesting to point out that the internal buf slice will grow at a rate of
// cap(buf)*2 + n.
// This means that if you've written 1MB into a buffer and then add 1 byte, your cap() will increase to 2097153 bytes.

func Test_VariableSizedBuffer(t *testing.T) {
	// 可变 size 的 Buffer
	// buf := &bytes.Buffer{}
	buf := new(bytes.Buffer)

	buf.WriteString("Hello World")
	bytes1 := buf.Bytes() // Bytes() return a slice, it's underlying
	fmt.Println("origin bytes1:", string(bytes1))
	str1 := buf.String()
	fmt.Println("origin str1:", str1)

	buf.Reset() // Reset 重置底层的写游标，但是底层的数据没有清除，寻址地址也没变。和上面 buf.Bytes() 返回的 slice 还是共享的底层数据。
	buf.WriteString("XXXX")
	bytes2 := buf.Bytes()
	fmt.Println("new bytes2", string(bytes2))

	fmt.Println("new bytes1:", string(bytes1))
	fmt.Println("new str1:", string(str1))

}

func Test_BufferRead(t *testing.T) {
	data := "Hello world\nHello again"

	//
	buf := bytes.NewBuffer([]byte(data))

	fmt.Println(buf.Bytes())

}

func Test_FixedSizedBuffer(t *testing.T) {
	// use a sized byte slice as a buffer
	buf := make([]byte, 1024)

	fmt.Println(buf)

}
