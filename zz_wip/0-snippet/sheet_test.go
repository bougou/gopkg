package snippet

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	xlsx "github.com/tealeg/xlsx/v3"
)

func Test_Sheet(t *testing.T) {

	filename := os.Args[1]

	f1, err := openXlsx(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("f1", f1)
	}

	sheet := f1.Sheets[0]

	for i := 0; i < sheet.MaxRow; i++ {
		row, _ := sheet.Row(i)
		fmt.Println(row)

	}
}

func openXlsx(filename string) (*xlsx.File, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return xlFile, nil
}

func openXlsx2(filename string) (*xlsx.File, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(len(fileBytes))
	fmt.Println(fileBytes)

	fileBytes = bytes.Trim(fileBytes, "\xef\xbb\xbf")
	fmt.Println(len(fileBytes))

	xlFile, err := xlsx.OpenBinary(fileBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return xlFile, nil
}
