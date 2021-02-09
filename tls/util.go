package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

var (
	ErrClientTLS = errors.New("client tls err")
)

func GetServerTLSConfig(caCertFile string) (t *tls.Config, err error) {
	if caCertFile == "" {
		return nil, nil
	}
	var ca *x509.CertPool
	data, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	ca = x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM(data); !ok {
		return nil, err
	}
	return &tls.Config{
		ClientCAs:  ca,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}, nil
}

func GetClientTLSConfig(insecure bool, caCertFile, CertFile, KeyFile string) (t *tls.Config, err error) {
	if caCertFile != "" || CertFile != "" || KeyFile != "" {
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
			return nil, ErrClientTLS
		}
		t.RootCAs = ca
	}
	if CertFile != "" && KeyFile != "" {
		crt, err := tls.LoadX509KeyPair(CertFile, KeyFile)
		if err != nil {
			return nil, err
		}
		t.Certificates = []tls.Certificate{crt}
	}
	return t, nil
}
