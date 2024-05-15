package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

func NewTLSConfigFromBytes(caBytes, certBytes, keyBytes []byte) (*tls.Config, error) {
	tlsCert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caBytes)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      pool,
	}

	return tlsConfig, nil
}

func NewTLSConfigFromFile(caFile, certFile, keyFile string) (*tls.Config, error) {
	caBytes, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certBytes, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	return NewTLSConfigFromBytes(caBytes, certBytes, keyBytes)
}

func NewHTTPSTransportFromFile(caFile, certFile, keyFile string) (*http.Transport, error) {
	tlsConfig, err := NewTLSConfigFromFile(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return transport, nil
}

func NewHTTPSTransportFromBytes(caBytes, certBytes, keyBytes []byte) (*http.Transport, error) {
	tlsConfig, err := NewTLSConfigFromBytes(caBytes, certBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return transport, nil
}
