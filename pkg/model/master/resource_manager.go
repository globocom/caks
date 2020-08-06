package master

import (
	"github.com/globocom/caks/pkg/model/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kuberesources "k8s.io/apimachinery/pkg/api/resource"
)

type ResourcesManager struct {}

func NewResourceSplitter() ResourcesManager {
	return ResourcesManager{}
}

func (r *ResourcesManager) sumDeploymentResources(deployment appsv1.Deployment)(*corev1.ResourceRequirements, error){
	var requirements  []corev1.ResourceRequirements

	for _, container := range deployment.Spec.Template.Spec.Containers {
		requirements = append(requirements,container.Resources)
	}

	return r.join(requirements...)
}

func (*ResourcesManager) join(resourcesRequirements ...corev1.ResourceRequirements)(*corev1.ResourceRequirements,error){

	sumCPURequest := 0
	sumCPULimit := 0
	sumMemoryRequest := 0
	sumMemoryLimit := 0

	for _, resource := range resourcesRequirements {
		memoryRequest, err := resources.ConvertToMebiBytes(resource.Requests.Memory().String())

		if err != nil {
			return nil, err
		}

		sumMemoryRequest += memoryRequest

		memoryLimit, err :=  resources.ConvertToMebiBytes(resource.Limits.Memory().String())

		if err != nil {
			return nil, err
		}

		sumMemoryLimit += memoryLimit

		cpuRequest, err := resources.ConvertToIntegerMiliCores(resource.Requests.Cpu().String())

		if err != nil {
			return nil, err
		}

		sumCPURequest += cpuRequest

		cpuLimit, err := resources.ConvertToIntegerMiliCores(resource.Limits.Cpu().String())

		if err != nil {
			return nil, err
		}

		sumCPULimit += cpuLimit
	}

	joinResource := &corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu": kuberesources.MustParse(
				resources.ConvertIntegerToStringMilicores(sumCPURequest),
			),
			"memory": kuberesources.MustParse(
				resources.ConvertIntegerToStringMebiBytes(sumMemoryRequest),
			),
		},
		Limits: corev1.ResourceList{
			"cpu": kuberesources.MustParse(
				resources.ConvertIntegerToStringMilicores(sumCPULimit),
			),
			"memory": kuberesources.MustParse(
				resources.ConvertIntegerToStringMebiBytes(sumMemoryLimit),
			),
		},
	}

	return joinResource, nil
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