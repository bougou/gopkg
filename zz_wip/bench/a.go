package bench

import (
	"os/exec"
	"path"

	archiver "github.com/mholt/archiver/v3"
	homedir "github.com/mitchellh/go-homedir"
)

func ByArchiver() {
	home, _ := homedir.Dir()
	s := path.Join(home, "Downloads", "tiny-imagenet-200.zip")
	d := path.Join(home, "Downloads", "test")
	archiver.Unarchive(s, d)
}

func ByUnzip() {
	home, _ := homedir.Dir()
	s := path.Join(home, "Downloads", "tiny-imagenet-200.zip")
	d := path.Join(home, "Downloads", "test")

	cmd := exec.Command("unzip", "-o", "-d", d, s)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
