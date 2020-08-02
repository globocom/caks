package master

import (
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type apiServer struct {
	image string
	applicationName string
	serviceClusterIPRange string
	advertiseAddress string
	admissionPlugins []string
	resourceRequirements corev1.ResourceRequirements
}

func (apiServer *apiServer) buildCommands()[]string{

	printAdmissionPlugins := func ()string{
		return strings.Join(apiServer.admissionPlugins,",")
	}

	return []string{
		apiServer.applicationName,
		printFlag("advertise-address",apiServer.advertiseAddress),
		printFlag("allow-privileged",true),
		printFlag("apiserver-count", 1),
		printFlag("audit-log-maxage",30),
		printFlag("audit-log-maxbackup",3),
		printFlag("audit-log-maxsize",100),
		printFlag("audit-log-path=","/var/log/audit.log"),
		printFlag("authorization-mode","Node,RBAC"),
		printFlag("bind-address","0.0.0.0"),
		printFlag("client-ca-file","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("enable-admission-plugins", printAdmissionPlugins()),
		printFlag("etcd-cafile","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("etcd-certfile","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("etcd-keyfile","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("etcd-servers","https://161.35.116.213:2379"),
		printFlag("event-ttl","1h"),
		printFlag("encryption-provider-config","/var/lib/kubernetes/encryption/encryption-config.yaml"),
		printFlag("kubelet-certificate-authority","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("kubelet-client-certificate","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("kubelet-client-key","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("kubelet-https",true),
		printFlag("runtime-config","api/all"),
		printFlag("service-account-key-file","/var/lib/kubernetes/service-account.pem"),
		printFlag("service-cluster-ip-range", apiServer.serviceClusterIPRange),
		printFlag("service-node-port-range","30000-32767"),
		printFlag("tls-cert-file","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("tls-private-key-file","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("v",2),
	}
}
