package master

import corev1 "k8s.io/api/core/v1"

type Scheduler struct {
	applicationName string
	image string
	resourceRequirements corev1.ResourceRequirements
}

func NewScheduler(resourceRequirements corev1.ResourceRequirements)Scheduler{
	return Scheduler{
		applicationName: "kube-scheduler",
		image: "rodrigoribeiro/globo-kube-scheduler",
		resourceRequirements: resourceRequirements,
	}
}