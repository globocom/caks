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
