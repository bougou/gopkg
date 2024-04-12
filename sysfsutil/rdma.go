package sysfsutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// getRdmaDevicePorts return the rdma device ports found from inside container
// each rdma device port is "{deviceName}:{portName}", like: mlx5_0:1, mlx5_1:1
func getRdmaDevicePorts() []string {
	results := []string{}

	filepath.WalkDir("/sys/class/infiniband", func(path string, d fs.DirEntry, _ error) error {
		fmt.Println("walk1", path)
		if path == "/sys/class/infiniband" {
			return nil
		}

		ok, _, err := isSymlink(path)
		if err != nil {
			// Just ignore
			return nil
		}

		if ok {
			device := d.Name()
			devicePortsPath := fmt.Sprintf("/sys/class/infiniband/%s/ports", device)

			dir, err := os.Open(devicePortsPath)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			defer dir.Close()

			items, err := dir.Readdir(0) // 0 means to return all items
			if err != nil {
				fmt.Println(err)
				return nil
			}

			for _, v := range items {
				port := v.Name()
				results = append(results, rdmaDevicePortToStr(device, port))
			}

		}

		return nil
	})

	return results
}

func isSymlink(absPath string) (ok bool, resolved string, err error) {
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		return false, "", err
	}

	if strings.HasPrefix(fileInfo.Mode().String(), "L") {
		resolved, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			return false, "", err
		}
		return true, resolved, nil
	}

	return false, "", err
}
func rdmaDevicePortToStr(rdmaDevice string, rdmaPort string) string {
	return fmt.Sprintf("%s:%s", rdmaDevice, rdmaPort)
}
