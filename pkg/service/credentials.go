package service

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
)

// Credentials - Credentials Service
type Credentials struct{}

//GenerateCertificate - Generate certificate
func (c *Credentials) GenerateCertificate(certificateAuthorityBytes []byte,
	clusterName string) ([]byte, []byte, error) {

	ca, err := c.readCertificate(certificateAuthorityBytes)

	if err != nil {
		return nil, nil, err
	}

	clusterCertificate := c.createCertificate(*ca, clusterName)

	privateKey, err := c.generateClientKey()

	if err != nil {
		return nil, nil, err
	}

	signedCertificate, err := c.assignCertificate(clusterCertificate, *ca, privateKey)

	if err != nil {
		return nil, nil, err
	}

	certificatePem := c.generateCertificatePem(signedCertificate)
	clientKeyPem := c.generatePrivateKeyPem(privateKey)

	return certificatePem, clientKeyPem, nil
}

func (c *Credentials) readCertificate(certificateBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certificateBytes)
	return x509.ParseCertificate(block.Bytes)
}

func (c *Credentials) createCertificate(ca x509.Certificate, clusterName string) x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization:  []string{clusterName},
			Country:       ca.Subject.Country,
			Province:      ca.Subject.Province,
			Locality:      ca.Subject.Locality,
			StreetAddress: ca.Subject.StreetAddress,
			PostalCode:    ca.Subject.PostalCode,
		},
	}
}

func (c *Credentials) generateClientKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func (c *Credentials) assignCertificate(certificate, ca x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	return x509.CreateCertificate(rand.Reader, &certificate, &ca, &privateKey.PublicKey, privateKey)
}

func (c *Credentials) generateCertificatePem(certificateByte []byte) []byte {
	return c.generatePem("CERTIFICATE", certificateByte)
}

func (c *Credentials) generatePrivateKeyPem(privateKey *rsa.PrivateKey) []byte {
	return c.generatePem("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
}

func (c *Credentials) generatePem(typ string, certificateByte []byte) []byte {

	var buffer bytes.Buffer

	pem.Encode(&buffer, &pem.Block{
		Type:  typ,
		Bytes: certificateByte,
	})

	return buffer.Bytes()
}
