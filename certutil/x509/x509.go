package x509

// Copied from: https://github.com/notaryproject/notary/blob/master/tuf/utils/x509.go

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// LoadCertFromPEM returns the first certificate found in a bunch of bytes or error
// if nothing is found. Taken from https://golang.org/src/crypto/x509/cert_pool.go#L85.
func LoadCertFromPEM(pemBytes []byte) (*x509.Certificate, error) {
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			return nil, errors.New("no certificates found in PEM data")
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}

		return cert, nil
	}

	return nil, errors.New("no certificates found in PEM data")
}

// NewCertificate returns an X509 Certificate following a template, given a Common Name and validity interval.
func NewCertificate(commonName string, startTime, endTime time.Time) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new certificate: %v", err)
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: startTime,
		NotAfter:  endTime,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}, nil
}

const (
	MinRSABitSize int = 2048
)

// ValidateCertificate returns an error if the certificate is not valid for notary
// Currently this is only ensuring the public key has a large enough modulus if RSA,
// using a non SHA1 signature algorithm, and an optional time expiry check
func ValidateCertificate(c *x509.Certificate, checkExpiry bool) error {
	if (c.NotBefore).After(c.NotAfter) {
		return fmt.Errorf("certificate validity window is invalid")
	}
	// Can't have SHA1 sig algorithm
	if c.SignatureAlgorithm == x509.SHA1WithRSA || c.SignatureAlgorithm == x509.DSAWithSHA1 || c.SignatureAlgorithm == x509.ECDSAWithSHA1 {
		return fmt.Errorf("certificate with CN %s uses invalid SHA1 signature algorithm", c.Subject.CommonName)
	}
	// If we have an RSA key, make sure it's long enough
	if c.PublicKeyAlgorithm == x509.RSA {
		rsaKey, ok := c.PublicKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("unable to parse RSA public key")
		}
		if rsaKey.N.BitLen() < MinRSABitSize {
			return fmt.Errorf("RSA bit length is too short")
		}
	}
	if checkExpiry {
		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)
		// Give one day leeway on creation "before" time, check "after" against today
		if (tomorrow).Before(c.NotBefore) || now.After(c.NotAfter) {
			return ErrCertExpired{CN: c.Subject.CommonName}
		}
	}
	return nil
}
