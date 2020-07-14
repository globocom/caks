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

