package sysfsutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const SysClassNetPath = "/sys/class/net"

func IsSymlink(absPath string) (ok bool, resolved string, err error) {
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

func GetPhysicalInterfaces() []string {
	infs := []string{}

	filepath.WalkDir("/sys/class/net", func(path string, d fs.DirEntry, _ error) error {
		if path == "/sys/class/net" {
			return nil
		}

		ok, resolved, err := IsSymlink(path)
		if err != nil {
			// Just ignore
			return nil
		}

		if ok {
			if !strings.Contains(resolved, "devices/virtual") {
				infs = append(infs, d.Name())
			}
		}
		return nil
	})

	return infs
}

func GetInterfaces() (pInfs []string, vInfs []string) {
	pInfs = []string{}
	vInfs = []string{}

	filepath.WalkDir("/sys/class/net", func(path string, d fs.DirEntry, _ error) error {
		if path == "/sys/class/net" {
			return nil
		}

		ok, resolved, err := IsSymlink(path)
		if err != nil {
			// Just ignore
			return nil
		}

		if ok {
			if strings.Contains(resolved, "devices/virtual") {
				vInfs = append(vInfs, d.Name())
			} else {
				pInfs = append(pInfs, d.Name())
			}
		}
		return nil
	})

	return pInfs, vInfs
}
