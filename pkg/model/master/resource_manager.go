package master

import (
	"github.com/globocom/caks/pkg/model/resources"
	corev1 "k8s.io/api/core/v1"
	kuberesources "k8s.io/apimachinery/pkg/api/resource"
)

type ResourcesManager struct {}

func NewResourceSplitter() ResourcesManager {
	return ResourcesManager{}
}

func (*ResourcesManager) split(controlPlaneResources corev1.ResourceRequirements,
	divisorStrategy func(res int)int)(*corev1.ResourceRequirements, error){

	res := corev1.ResourceRequirements{}
	requests := controlPlaneResources.Requests
	limits := controlPlaneResources.Limits

	cpuRequestsValue, err := resources.ConvertToIntegerMiliCores(requests.Cpu().String())

	if err != nil {
		return nil, err
	}

	cpuLimitValue, 	err := resources.ConvertToIntegerMiliCores(limits.Cpu().String())

	if err != nil {
		return nil, err
	}

	memoryRequestsValue, err := resources.ConvertToMebiBytes(requests.Memory().String())

	if err != nil {
		return nil, err
	}

	memoryLimitValue, err := resources.ConvertToMebiBytes(limits.Memory().String())

	if err != nil {
		return nil, err
	}

	res.Requests = corev1.ResourceList{
		"cpu": kuberesources.MustParse(
			resources.ConvertIntegerToStringMilicores(divisorStrategy(cpuRequestsValue))),
		"memory": kuberesources.MustParse(
			resources.ConvertIntegerToStringMebiBytes(divisorStrategy(memoryRequestsValue))),
	}

	res.Limits = corev1.ResourceList{
		"cpu": kuberesources.MustParse(
			resources.ConvertIntegerToStringMilicores(divisorStrategy(cpuLimitValue))),
		"memory": kuberesources.MustParse(
			resources.ConvertIntegerToStringMebiBytes(divisorStrategy(memoryLimitValue))),
	}

	return &res, nil
}