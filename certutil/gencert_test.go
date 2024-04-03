package certutil

import (
	"fmt"
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func Test_Gen(t *testing.T) {
	home, _ := homedir.Dir()
	outputDir := home

	validDays := 999 * 365
	ca, err := NewCA("ca", validDays)
	if err != nil {
		t.Error("NewCA failed: ", err)
	}
	if err != ca.Dump(outputDir) {
		t.Error(err)
	}

	server := NewKeyCert("127.0.0.1")
	if err := server.GenKey(); err != nil {
		t.Error(err)
	}
	if server.SignedByCA(ca, validDays); err != nil {
		t.Error(err)
	}
	if err != server.Dump(outputDir) {
		t.Error(err)
	}
}

func Test_LoadCertificateDERFile(t *testing.T) {
	home, _ := homedir.Dir()

	cert, err := LoadCertficateDERFile(path.Join(home, "AAA-Certificate.cer"))
	if err != nil {
		t.Error(err)
	}
	kc := &KeyCert{cert: cert}
	b, err := kc.GetCertPEMBytes()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
