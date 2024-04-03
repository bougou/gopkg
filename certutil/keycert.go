package certutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"

	valid "github.com/asaskevich/govalidator"
)

// KeyCert holds rsa.PrivateKey and x509.Certificate together.
type KeyCert struct {
	cert *x509.Certificate
	key  *rsa.PrivateKey
}

// NewKeyCert creates an emtpy KeyCert.
// You must use its GenKey method to fill its private key
// and use a CA to sign it to fill its cert.
func NewKeyCert(commonName string, options ...DNOption) *KeyCert {
	kc := &KeyCert{
		cert: new(x509.Certificate),
		key:  new(rsa.PrivateKey),
	}

	kc.cert.Subject.CommonName = commonName
	kc.cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	kc.cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	kc.cert.BasicConstraintsValid = true
	if valid.IsIP(commonName) {
		ip := net.ParseIP(commonName)
		kc.cert.IPAddresses = append(kc.cert.IPAddresses, ip)
	}
	if valid.IsDNSName(commonName) {
		kc.cert.DNSNames = append(kc.cert.DNSNames, commonName)
	}

	for _, option := range options {
		option(kc)
	}
	return kc
}

func LoadKeyCertPEMFile(keyFile string, certFile string) (*KeyCert, error) {
	keyPEMBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	certPEMBytes, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	return LoadKeyCertPEMBytes(keyPEMBytes, certPEMBytes)
}

// LoadKeyCertPEMBytes loads key and cert from pem bytes, and return a fully filled KeyCert.
func LoadKeyCertPEMBytes(keyPEMBytes []byte, certPEMBytes []byte) (*KeyCert, error) {
	key, err := LoadKeyPEMBytes(keyPEMBytes)
	if err != nil {
		return nil, err
	}

	cert, err := LoadCertificatePEMBytes(certPEMBytes)
	if err != nil {
		return nil, err
	}

	kc := &KeyCert{
		key:  key,
		cert: cert,
	}
	return kc, nil
}

func LoadKeyPEMBytes(keyPEMBytes []byte) (*rsa.PrivateKey, error) {
	keyPEMBlock, _ := pem.Decode(keyPEMBytes)
	if keyPEMBlock == nil {
		return nil, errors.New("not valid key")
	}
	key, err := x509.ParsePKCS1PrivateKey(keyPEMBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("ParsePKCS1PrivateKey failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return key, nil
}

func LoadKeyPEMFile(keyFile string) (*rsa.PrivateKey, error) {
	keyPEMBytes, err := os.ReadFile(keyFile)
	if err != nil {
		msg := fmt.Sprintf("Read key file failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return LoadKeyPEMBytes(keyPEMBytes)
}

func LoadCertificatePEMBytes(certPEMBytes []byte) (*x509.Certificate, error) {
	certPEMBlock, _ := pem.Decode(certPEMBytes)
	if certPEMBlock == nil {
		return nil, errors.New("not valid cert")
	}
	cert, err := x509.ParseCertificate(certPEMBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("ParseCertificate file failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return cert, nil
}

func LoadCertificatePEMFile(certFile string) (*x509.Certificate, error) {
	certPEMBytes, err := os.ReadFile(certFile)
	if err != nil {
		msg := fmt.Sprintf("Read Cert file failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return LoadCertificatePEMBytes(certPEMBytes)
}

func LoadCertficateDERFile(certFile string) (*x509.Certificate, error) {
	certDERBytes, err := os.ReadFile(certFile)
	if err != nil {
		msg := fmt.Sprintf("Read Cert file failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return LoadCertficateDERBytes(certDERBytes)
}

func LoadCertficateDERBytes(certDERBytes []byte) (*x509.Certificate, error) {
	cert, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		msg := fmt.Sprintf("ParseCertificate file failed, err: %s", err)
		return nil, errors.New(msg)
	}
	return cert, nil
}
