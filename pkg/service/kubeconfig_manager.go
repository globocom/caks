package service

import (
	"bytes"
	"encoding/base64"
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

//KubeConfigBuilder - kubeconfig generator
type KubeConfigBuilder struct {
}

//GenerateKubeConfig - generate kubernetes configuration
func (kcm *KubeConfigBuilder) GenerateKubeConfig(credential Credential, server, user string) (*string, error) {

	certificateBase64 := kcm.generateBase64(credential.certificatePEM)
	clientKeyBase64 := kcm.generateBase64(credential.clientKeyPEM)
	caBase64 := kcm.generateBase64(credential.certificateAuthority)

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
		CA:                caBase64,
		Server:            server,
		ClusterName:       credential.clusterName,
		User:              user,
		ClientCertificate: certificateBase64,
		ClientKey:         clientKeyBase64,
	}

	err = template.Execute(&kubeConfigBuffer, kubeConfigData)

	if err != nil {
		return nil, err
	}

	kubeConfig := kubeConfigBuffer.String()

	return &kubeConfig, nil
}

func (kcm *KubeConfigBuilder) generateBase64(certificate []byte) string {
	return base64.RawStdEncoding.EncodeToString(certificate)
}
