package copy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// Whether the srcDir string ends with '/' determines the copy behaviour.
// It des not matter whether the dstDir string ends with '/'.
//
// Both the srcDir and the dstDir must be absolute paths which starts with '/'.
//
// If the srcDir ends in a '/', the contents of the directory are copied rather than the directory itself.
//
//	eg 1:
//	srcDir: /path/to/src1/     {file1.txt,...}
//	dstDir: /path/to/dst1
//	result: /path/to/dst1/{file1.txt,...}
//
//	eg 2:
//	srcDir: /path/to/src1     {file1.txt,...}
//	dstDir: /path/to/dst1
//	result: /path/to/dst1/src1/{file1.txt,...}
func CopyDir(srcDir string, dstDir string) error {
	if !strings.HasPrefix(srcDir, "/") || !strings.HasPrefix(dstDir, "/") {
		return fmt.Errorf("both srcDir and dstDir must be absolute dir path")
	}

	if path.Join(srcDir) == path.Join(dstDir) {
		// ignore
		return nil
	}

	// if srcDir ends with "/", then srcParent == srcDir (with end / removed)
	// if srcDir not ends with "/", then srcParent == parent dir of srcDir
	srcParent := filepath.Dir(srcDir)

	err := filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// srcfile relative path
		relPath, err := filepath.Rel(srcParent, srcPath)
		if err != nil {
			return err
		}

		stat, ok := f.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("unable to get raw syscall.Stat_t data for %s", srcPath)
		}

		dstPath := filepath.Join(dstDir, relPath)

		switch mode := f.Mode(); {
		case mode.IsDir():
			if err := os.MkdirAll(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
				return err
			}

		case mode.IsRegular():
			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(dstPath, input, f.Mode())
			if err != nil {
				return err
			}

		case mode&os.ModeSymlink != 0:
			link, err := os.Readlink(srcPath)
			if err != nil {
				return err
			}

			if err := os.Symlink(link, dstPath); err != nil {
				return err
			}

		case mode&os.ModeNamedPipe != 0:
			fallthrough

		case mode&os.ModeSocket != 0:
			if err := unix.Mkfifo(dstPath, uint32(stat.Mode)); err != nil {
				return err
			}

		case mode&os.ModeDevice != 0:
			if err := unix.Mknod(dstPath, uint32(stat.Mode), int(stat.Rdev)); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown file type (%d / %s) for %s", f.Mode(), f.Mode().String(), srcPath)
		}

		return nil
	})

	return err
}

func CopyFile(srcFile, dstFile string) error {
	f, err := os.Lstat(srcFile)
	if err != nil {
		return err
	}

	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dstFile, input, f.Mode())
	if err != nil {
		return err
	}

	return nil
}
