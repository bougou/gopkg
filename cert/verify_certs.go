package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func VerifyCerts(certPEM []byte, rootPEMs ...[]byte) (bool, error) {

	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	for _, rootPEM := range rootPEMs {
		ok := roots.AppendCertsFromPEM([]byte(rootPEM))
		if !ok {
			return false, errors.New("failed to parse root certificate")
		}
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return false, errors.New("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, errors.New("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		// DNSName: "mail.google.com",
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return false, errors.New("failed to verify certificate: " + err.Error())
	}

	return true, nil
}
