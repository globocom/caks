package service

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
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

	clusterCertificate := kcm.createCertificate(clusterName)

	privateKey, err := kcm.generateClientKey()

	if err != nil {
		return nil, err
	}

	signedCertificate, err := kcm.assignCertificate(clusterCertificate, *ca, privateKey)

	if err != nil {
		return nil, err
	}

	certificatePem := kcm.generateBase64Pem(signedCertificate)
	clientKeyPem := kcm.generateBase64Pem(x509.MarshalPKCS1PrivateKey(privateKey))
	caPem := kcm.generateBase64Pem(certificateAuthorityBytes)

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
	return nil, nil
}

func (kcm *KubeConfigManager) createCertificate(clusterName string) x509.Certificate {
	return x509.Certificate{}
}

func (kcm *KubeConfigManager) generateClientKey() (*rsa.PrivateKey, error) {
	return nil, nil
}

func (kcm *KubeConfigManager) assignCertificate(certificate, ca x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	return nil, nil
}

func (kcm *KubeConfigManager) generateBase64Pem(certificateByte []byte) string {
	return ""
}
