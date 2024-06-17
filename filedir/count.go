package filedir

import (
	"fmt"
	"os"
	"path/filepath"
)

func CountItemsOfDir(dirPath string) (int, error) {
	file, err := os.Open(dirPath)
	if err != nil {
		return 0, fmt.Errorf("open dir %s failed, err: %s", dirPath, err)
	}

	files, err := file.Readdirnames(-1)
	if err != nil {
		return 0, err
	}

	return len(files), nil
}

func CountItemsOfDirRecusively(dirPath string) (int, error) {
	count := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Increment the count for each item
		count++
		return nil
	})

	if err != nil {
		return 0, err
	}

	// Subtract 1 from the count to exclude the root directory itself
	count--
	return count, nil
}
