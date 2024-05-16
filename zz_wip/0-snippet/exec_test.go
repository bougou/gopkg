package snippet

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func exec1() {
	cmd := exec.Command("ls", "-l", "-G")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func exec2() {
	// On Linux --color=always force color
	// On Mac -G force color
	cmd1 := exec.Command("ls", "-G", "-l", ".")
	cmd1.Env = os.Environ()
	var o, e bytes.Buffer
	cmd1.Stdout = &o
	cmd1.Stderr = &e
	fmt.Println(cmd1)
	cmd1.Run()
	os.WriteFile("/tmp/test.txt", o.Bytes(), 0600)

}
