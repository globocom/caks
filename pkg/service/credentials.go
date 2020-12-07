package service

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
)

// Credentials - Credentials Service
type Credentials struct{}

//GenerateCertificate - Generate certificate
func (c *Credentials) GenerateCertificate(certificateAuthorityBytes []byte,
	clusterName string) (certificate []byte, privateKey []byte) {

	return nil, nil
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
