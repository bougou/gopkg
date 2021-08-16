package file

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// WriteToCSV write data to csv file
func WriteToCSV(data *[][]string, filename string) error {
	fileDir := filepath.Dir(filename)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		msg := fmt.Sprintf("create dir (%s) failed, err: %s", fileDir, err)
		return errors.New(msg)
	}

	file, err := os.Create(filename)
	if err != nil {
		msg := fmt.Sprintf("create file (%s) failed, err: %s", filename, err)
		return errors.New(msg)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i, value := range *data {
		err := writer.Write(value)
		if err != nil {
			msg := fmt.Sprintf("write the %d row failed, err: %s", i+1, err)
			return errors.New(msg)
		}
	}

	return nil
}
