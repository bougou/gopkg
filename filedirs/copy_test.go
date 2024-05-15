package filedirs

import (
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func Test_c(t *testing.T) {
	home, _ := homedir.Dir()

	srcDir := path.Join(home, "srctest")
	dstDir := path.Join(home, "dsttest")
	if err := CopyDir(srcDir, dstDir); err != nil {
		t.Error(err)
	}
}
