package master

import (
	corev1 "k8s.io/api/core/v1"
)

type ControllerManager struct {
	applicationName string
	image string
	clusterCIDRS string
	clusterName string
	serviceClusterIPRange string
	resourceRequirements corev1.ResourceRequirements
}

func (*ControllerManager) buildVolumeMounts()[]corev1.VolumeMount{
	return []corev1.VolumeMount{
		{Name: "kubernetes", MountPath: "/var/lib/kubernetes", ReadOnly: true},
		{Name: "ca", MountPath: "/var/lib/kubernetes/ca", ReadOnly: true},
	}
}

func (controllerManager *ControllerManager) buildCommands()[]string{
	return []string{
		controllerManager.applicationName,
		printFlag("address", "0.0.0.0"),
		printFlag("cluster-cidr",controllerManager.clusterCIDRS),
		printFlag("allocate-node-cidrs", true),
		printFlag("cluster-name",controllerManager.clusterName),
		printFlag("cluster-signing-cert-file","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("cluster-signing-key-file", "/var/lib/kubernetes/ca/ca-key.pem"),
		printFlag("kubeconfig","/var/lib/kubernetes/kube-controller-manager.kubeconfig"),
		printFlag("leader-elect", true),
		printFlag("root-ca-file","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("service-account-private-key-file","/var/lib/kubernetes/service-account-key.pem"),
		printFlag("service-cluster-ip-range",controllerManager.serviceClusterIPRange),
		printFlag("use-service-account-credentials", true),
		printFlag("v", 2),
	}
}
