package master

import (
	"github.com/globocom/caks/pkg/apis/cacks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type Master struct{
	settings v1alpha1.ControlPlaneMaster
	namespacedName types.NamespacedName
	apiServer apiServer
	scheduler Scheduler
	controllerManager ControllerManager
	resourceManager ResourcesManager
}

func (master *Master) buildPod()corev1.PodTemplateSpec{

	return corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{
			Namespace: master.namespacedName.Namespace,
			Labels: master.buildPodLabels(),
		},
		Spec: corev1.PodSpec{
			Volumes: master.buildVolumes(),
			Containers: []corev1.Container{
				master.apiServer.BuildContainer(),
				master.scheduler.BuilderContainer(),
				master.controllerManager.BuilderContainer(),
			},
		},
	}
}

func (master *Master) buildVolumes()[]corev1.Volume{

	return []corev1.Volume{
		master.buildSecretVolume("ca", "ca-certs"),
		master.buildSecretVolume("kubernetes", master.settings.MasterSecretName),
		master.buildSecretVolume("encryption", master.settings.EncryptionSecretName),
	}
}

func (master *Master) buildPodLabels()map[string]string{
	return map[string]string{
		"app":"master",
		"cluster": master.namespacedName.Name,
		"tier": "control-plane",
	}
}