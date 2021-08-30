package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"
	"time"
)

const (
	KeyFileSuffix  = ".key"
	CertFileSuffix = ".crt"
	PemFileSuffix  = ".pem"
)

func randomSerialNumber() (*big.Int, error) {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(100), nil).Sub(max, big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		msg := fmt.Sprintf("generating serial number failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return n, nil
}

// GenKey fill kc.key and kc.keyPEMBytes
func (kc *KeyCert) GenKey() error {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	kc.key = key
	return nil
}

func NewCA(caName string, validDays int, options ...DNOption) (*KeyCert, error) {
	ca := NewKeyCert(caName, options...)
	if err := ca.GenKey(); err != nil {
		msg := fmt.Sprintf("GenKey failed, err: %s", err)
		return nil, errors.New(msg)
	}
	ca.cert.IsCA = true

	sn, err := randomSerialNumber()
	if err != nil {
		msg := fmt.Sprintf("gen randomSerialNumber failed, err: %s", err)
		return nil, errors.New(msg)
	}
	ca.cert.SerialNumber = sn

	ca.cert.NotBefore = time.Now()
	ca.cert.NotAfter = time.Now().AddDate(0, 0, validDays)
	ca.cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ca.cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	ca.cert.BasicConstraintsValid = true

	caCertBytes, err := x509.CreateCertificate(rand.Reader, ca.cert, ca.cert, &ca.key.PublicKey, ca.key)
	if err != nil {
		msg := fmt.Sprintf("CreateCertfiicate failed, err: %s", err)
		return nil, errors.New(msg)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBytes,
	}
	caCert, err := x509.ParseCertificate(caCertPEMBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("ParseCertificate failed, err: %s", err)
		return nil, errors.New(msg)
	}

	ca.cert = caCert
	return ca, nil
}

func (kc *KeyCert) SignedByCA(ca *KeyCert, validDays int) error {
	sn, err := randomSerialNumber()
	if err != nil {
		msg := fmt.Sprintf("gen randomSerialNumber failed, err: %s", err)
		return errors.New(msg)
	}
	kc.cert.SerialNumber = sn
	kc.cert.NotBefore = time.Now()
	kc.cert.NotAfter = time.Now().AddDate(0, 0, validDays)

	certBytes, err := x509.CreateCertificate(rand.Reader, kc.cert, ca.cert, &kc.key.PublicKey, ca.key)
	if err != nil {
		msg := fmt.Sprintf("CreateCertfiicate failed, err: %s", err)
		return errors.New(msg)
	}
	certPEM := new(bytes.Buffer)
	certPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	pem.Encode(certPEM, certPEMBlock)
	cert, err := x509.ParseCertificate(certPEMBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("ParseCertificate failed, err: %s", err)
		return errors.New(msg)
	}
	kc.cert = cert
	return nil
}

func (ca *KeyCert) Sign(kc *KeyCert, validDays int) error {
	if !ca.cert.IsCA {
		return errors.New("not a CA, only CA can sign new certs")
	}

	if kc.key == nil {
		return errors.New("key must be initialized, use GenKey() method to initialize key")
	}
	kc.SignedByCA(ca, validDays)
	return nil
}

func (kc *KeyCert) Dump(outputDir string) error {
	keyPEMBytes, err := kc.GetKeyPEMBytes()
	if err != nil {
		return err
	}
	certPEMBytes, err := kc.GetCertPEMBytes()
	if err != nil {
		return err
	}
	commonName := kc.cert.Subject.CommonName
	if commonName == "" {
		return errors.New("not valid kc, the subject common name of cert must not be empty")
	}
	commonName = strings.ReplaceAll(commonName, " ", "-")

	keyFile := path.Join(outputDir, commonName+KeyFileSuffix)
	certFile := path.Join(outputDir, commonName+CertFileSuffix)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(keyFile, keyPEMBytes, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(certFile, certPEMBytes, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (kc *KeyCert) GetKeyPEMBytes() ([]byte, error) {
	keyPEM := new(bytes.Buffer)
	err := pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(kc.key),
	})
	if err != nil {
		return nil, err
	}
	return keyPEM.Bytes(), nil
}

func (kc *KeyCert) GetCertPEMBytes() ([]byte, error) {
	certPEM := new(bytes.Buffer)
	certPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: kc.cert.Raw,
	}
	err := pem.Encode(certPEM, certPEMBlock)
	if err != nil {
		return nil, err
	}

	return certPEM.Bytes(), nil
}
