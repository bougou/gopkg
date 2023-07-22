package bit

import (
	"fmt"
	"strconv"
)

func FormatByteSlice(data []byte, width int) string {
	res := ""
	for k, v := range data {
		res += fmt.Sprintf("%02x", v)

		if k%width == width-1 {
			res += "\n"
		} else {
			res += " "
		}
	}
	return res
}

// w 条带宽度
// BitPrint(65536, 8)
// 00000001 00000000 00000000
func BitPrint(i int64, w int) {
	s := []byte(strconv.FormatInt(i, 2))

	// 需要补齐的位数
	f := w - len(s)%w
	if f != 0 {
		o := make([]byte, f)
		for i := range o {
			o[i] = '0'
		}
		s = append(o, s...)
	}

	for i, v := range s {
		// 开始显示下一个条带时，打印一个间隔
		if i%w == 0 && i != 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%c", v)
	}
	fmt.Println()
}
