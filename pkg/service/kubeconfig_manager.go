package service

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"text/template"
)

const kubeConfigTemplate = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{.CA}}
    server: {{.Server}}
  name: {{.ClusterName}}
contexts:
- context:
    cluster: {{.ClusterName}}
    user: "system:"{{.User}}
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: "system:"{{.User}}
  user:
    client-certificate-data: {{.ClientCertificate}}
    client-key-data: {{.ClientKey}}
`

//KubeConfigManager - kubeconfig generator
type KubeConfigManager struct {
}

func (kcm *KubeConfigManager) GenerateKubeConfig(certificateAuthorityBytes []byte, server,
	clusterName, user string) (*string, error) {

	ca, err := kcm.readCertificate(certificateAuthorityBytes)

	if err != nil {
		return nil, err
	}

	clusterCertificate := kcm.createCertificate(*ca, clusterName)

	privateKey, err := kcm.generateClientKey()

	if err != nil {
		return nil, err
	}

	signedCertificate, err := kcm.assignCertificate(clusterCertificate, *ca, privateKey)

	if err != nil {
		return nil, err
	}

	certificatePem := kcm.generateCertificateBase64Pem(signedCertificate)
	clientKeyPem := kcm.generatePrivateKeyBase64Pem(privateKey)
	caPem := kcm.generateCertificateBase64Pem(certificateAuthorityBytes)

	template, err := template.New("kubeconfig").Parse(kubeConfigTemplate)

	if err != nil {
		return nil, err
	}

	var kubeConfigBuffer bytes.Buffer

	kubeConfigData := struct {
		CA                string
		Server            string
		ClusterName       string
		User              string
		ClientCertificate string
		ClientKey         string
	}{
		CA:                caPem,
		Server:            server,
		ClusterName:       clusterName,
		User:              user,
		ClientCertificate: certificatePem,
		ClientKey:         clientKeyPem,
	}

	err = template.Execute(&kubeConfigBuffer, kubeConfigData)

	if err != nil {
		return nil, err
	}

	kubeConfig := kubeConfigBuffer.String()

	return &kubeConfig, nil
}

func (kcm *KubeConfigManager) readCertificate(certificateBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certificateBytes)
	return x509.ParseCertificate(block.Bytes)
}

func (kcm *KubeConfigManager) createCertificate(ca x509.Certificate, clusterName string) x509.Certificate {
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

func (kcm *KubeConfigManager) generateClientKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func (kcm *KubeConfigManager) assignCertificate(certificate, ca x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	return x509.CreateCertificate(rand.Reader, &certificate, &ca, privateKey.PublicKey, privateKey)
}

func (kcm *KubeConfigManager) generateCertificateBase64Pem(certificateByte []byte) string {
	return kcm.generateBase64Pem("CERTIFICATE", certificateByte)
}

func (kcm *KubeConfigManager) generatePrivateKeyBase64Pem(privateKey *rsa.PrivateKey) string {
	return kcm.generateBase64Pem("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
}

func (kcm *KubeConfigManager) generateBase64Pem(typ string, certificateByte []byte) string {

	var buffer bytes.Buffer

	pem.Encode(&buffer, &pem.Block{
		Type:  typ,
		Bytes: certificateByte,
	})

	return base64.RawStdEncoding.EncodeToString(buffer.Bytes())
}
