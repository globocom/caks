package master

import corev1 "k8s.io/api/core/v1"

type apiServer struct {
	image string
	applicationName string
	serviceClusterIPRange string
	advertiseAddress string
	admissionPlugins []string
	resourceRequirements corev1.ResourceRequirements
}
