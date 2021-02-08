package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

var (
	ErrTLS = errors.New("tls err")
)

func GetTLSConfig(insecure bool, caCertFile, clientCertFile, clientKeyFile string) (t *tls.Config, err error) {
	if caCertFile != "" || clientCertFile != "" || clientKeyFile != "" {
		t = &tls.Config{
			InsecureSkipVerify: insecure,
		}
	} else {
		return nil, nil
	}
	if caCertFile != "" {
		var ca *x509.CertPool
		data, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		ca = x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM(data); !ok {
			return nil, ErrTLS
		}
		t.RootCAs = ca
	}
	if clientCertFile != "" && clientKeyFile != "" {
		crt, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, err
		}
		t.Certificates = []tls.Certificate{crt}
	}
	return t, nil
}
