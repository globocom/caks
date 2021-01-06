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

// Credential - cluster credential
type Credential struct {
	certificateAuthority []byte
	certificatePEM       []byte
	clientKeyPEM         []byte
	clusterName          string
}

// CredentialsBuilder - Credentials Service
type CredentialsBuilder struct{}

//Build - Build credential
func (c *CredentialsBuilder) Build(certificateAuthority []byte,
	clusterName string) (*Credential, error) {

	ca, err := c.readCertificate(certificateAuthority)

	if err != nil {
		return nil, err
	}

	clusterCertificate := c.createCertificate(*ca, clusterName)

	privateKey, err := c.generateClientKey()

	if err != nil {
		return nil, err
	}

	signedCertificate, err := c.assignCertificate(clusterCertificate, *ca, privateKey)

	if err != nil {
		return nil, err
	}

	certificatePem := c.generateCertificatePem(signedCertificate)
	clientKeyPem := c.generatePrivateKeyPem(privateKey)

	return &Credential{
		clientKeyPEM:         clientKeyPem,
		certificatePEM:       certificatePem,
		certificateAuthority: certificateAuthority,
		clusterName:          clusterName,
	}, nil
}

func (c *CredentialsBuilder) readCertificate(certificateBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certificateBytes)
	return x509.ParseCertificate(block.Bytes)
}

func (c *CredentialsBuilder) createCertificate(ca x509.Certificate, clusterName string) x509.Certificate {
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

func (c *CredentialsBuilder) generateClientKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func (c *CredentialsBuilder) assignCertificate(certificate, ca x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	return x509.CreateCertificate(rand.Reader, &certificate, &ca, &privateKey.PublicKey, privateKey)
}

func (c *CredentialsBuilder) generateCertificatePem(certificateByte []byte) []byte {
	return c.generatePem("CERTIFICATE", certificateByte)
}

func (c *CredentialsBuilder) generatePrivateKeyPem(privateKey *rsa.PrivateKey) []byte {
	return c.generatePem("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
}

func (c *CredentialsBuilder) generatePem(typ string, certificateByte []byte) []byte {

	var buffer bytes.Buffer

	pem.Encode(&buffer, &pem.Block{
		Type:  typ,
		Bytes: certificateByte,
	})

	return buffer.Bytes()
}
