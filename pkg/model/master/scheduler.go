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

func (scheduler *Scheduler) BuilderContainer()corev1.Container{
	return corev1.Container{
		Name: scheduler.applicationName,
		Image: scheduler.image,
		Command: scheduler.buildCommands(),
		VolumeMounts: scheduler.buildVolumeMounts(),
		Resources: scheduler.resourceRequirements,
	}
}

func (*Scheduler) buildVolumeMounts()[]corev1.VolumeMount{
	return []corev1.VolumeMount{
		{Name: "kubernetes", MountPath: "/var/lib/kubernetes", ReadOnly: true},
	}
}

func (scheduler *Scheduler) buildCommands()[]string{
	return []string{
		scheduler.applicationName,
		printFlag("leader-elect", true),
		printFlag("kubeconfig","/var/lib/kubernetes/kube-scheduler.kubeconfig"),
		printFlag("v", 2),
	}
}
