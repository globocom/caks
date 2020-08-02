package master

import (
	"github.com/globocom/caks/pkg/apis/cacks/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)


type Master struct{
	settings v1alpha1.ControlPlaneMaster
	namespacedName types.NamespacedName
	apiServer apiServer
	scheduler Scheduler
	controllerManager ControllerManager
	resourceManager ResourcesManager
}