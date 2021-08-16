package file

import (
	"errors"
	"fmt"

	"github.com/tealeg/xlsx"
)

var defaultXLSXSheetName string = "default"

// WriteToXLSX write data to xlsx file
func WriteToXLSX(data *[][]string, filename string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(defaultXLSXSheetName)
	if err != nil {
		msg := fmt.Sprintf("create sheet failed, err: %s", err)
		return errors.New(msg)
	}

	for _, rowData := range *data {
		row := sheet.AddRow()
		// row.SetHeightCM(0.5)

		for _, cellData := range rowData {
			cell := row.AddCell()
			cell.Value = cellData
		}
	}

	if err := file.Save(filename); err != nil {
		msg := fmt.Sprintf("save file (%s) failed, err: %s", filename, err)
		return errors.New(msg)
	}

	return nil
}
