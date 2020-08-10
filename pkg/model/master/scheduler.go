package master

import corev1 "k8s.io/api/core/v1"

type Scheduler struct {
	applicationName string
	image string
	resourceRequirements corev1.ResourceRequirements
}