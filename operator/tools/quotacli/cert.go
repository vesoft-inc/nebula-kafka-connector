package main

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	"k8s.io/klog/v2"
)

const (
	rsaKeySize                  = 2048
	RSAPrivateKeyBlockType      = "RSA PRIVATE KEY"
	CertificateRequestBlockType = "CERTIFICATE REQUEST"
)

type FileProjection struct {
	Data []byte
	Mode int32
}

func WriteCertsToDir(dir string, user string, key, csr []byte) error {
	keyName := getKeyFileName(user)
	csrName := getCsrFileName(user)
	if err := prepareToWrite(dir, keyName, csrName); err != nil {
		return err
	}
	return writePayloadToDir(certToProjectionMap(user, key, csr), dir)
}

func prepareToWrite(dir string, filenames ...string) error {
	_, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		klog.Infof("cert directory doesn't exist, creating directory <%s>", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("can't create dir: %v", dir)
		}
	case err != nil:
		return err
	}

	for _, f := range filenames {
		abspath := path.Join(dir, f)
		_, err := os.Stat(abspath)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			klog.Errorf("unable to stat file %s: %v", abspath, err)
		}
		err = os.Remove(abspath)
		if err != nil {
			klog.Errorf("unable to remove old file %s: %v", abspath, err)
		}
	}
	return nil
}

func writePayloadToDir(payload map[string]FileProjection, dir string) error {
	for fileName, fileProjection := range payload {
		content := fileProjection.Data
		mode := os.FileMode(fileProjection.Mode)
		fullPath := path.Join(dir, fileName)

		err := os.WriteFile(fullPath, content, mode)
		if err != nil {
			klog.Errorf("unable to write file %s: %v", fullPath, err)
			return err
		}

		err = os.Chmod(fullPath, mode)
		if err != nil {
			klog.Errorf("unable to chmod file %s: %v", fullPath, err)
		}
	}

	return nil
}

func certToProjectionMap(user string, key, csr []byte) map[string]FileProjection {
	keyName := getKeyFileName(user)
	csrName := getCsrFileName(user)
	return map[string]FileProjection{
		keyName: {
			Data: key,
			Mode: 0644,
		},
		csrName: {
			Data: csr,
			Mode: 0644,
		},
	}
}

func GenerateCert(commonName string) ([]byte, []byte, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the private key: %v", err)
	}

	csr, err := MakeCSR(key, commonName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the certificate signing request: %v", err)
	}

	return EncodePrivateKeyPEM(key), csr, nil
}

func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
}

func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

func MakeCSR(privateKey interface{}, commonName string) (csr []byte, err error) {
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
		},
	}

	return MakeCSRFromTemplate(privateKey, template)
}

func MakeCSRFromTemplate(privateKey interface{}, template *x509.CertificateRequest) ([]byte, error) {
	t := *template
	t.SignatureAlgorithm = sigType(privateKey)

	csrDER, err := x509.CreateCertificateRequest(cryptorand.Reader, &t, privateKey)
	if err != nil {
		return nil, err
	}

	csrPemBlock := &pem.Block{
		Type:  CertificateRequestBlockType,
		Bytes: csrDER,
	}

	return pem.EncodeToMemory(csrPemBlock), nil
}

func sigType(privateKey interface{}) x509.SignatureAlgorithm {
	// Customize the signature for RSA keys, depending on the key size
	if privateKey, ok := privateKey.(*rsa.PrivateKey); ok {
		keySize := privateKey.N.BitLen()
		switch {
		case keySize >= 4096:
			return x509.SHA512WithRSA
		case keySize >= 3072:
			return x509.SHA384WithRSA
		default:
			return x509.SHA256WithRSA
		}
	}
	return x509.UnknownSignatureAlgorithm
}
