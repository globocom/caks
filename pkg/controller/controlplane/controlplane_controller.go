package controlplane

import (
	"context"
	"github.com/globocom/caks/pkg/model/master"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/util/slice"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	caksv1alpha1 "github.com/globocom/caks/pkg/apis/caks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_controlplane")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileControlPlane{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("controlplane-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ControlPlane
	err = c.Watch(&source.Kind{Type: &caksv1alpha1.ControlPlane{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Deployment and requeue the owner ControlPlane
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{

	}, predicate.GenerationChangedPredicate{Funcs: predicate.Funcs{DeleteFunc: func(e event.DeleteEvent) bool{
		if _, ok := e.Meta.GetLabels()["tier"]; ok {
			return true
		}
		return false
	}}})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner ControlPlane
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &caksv1alpha1.ControlPlane{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileControlPlane{}

// ReconcileControlPlane reconciles a ControlPlane object
type ReconcileControlPlane struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileControlPlane) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ControlPlane")

	instance := &caksv1alpha1.ControlPlane{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err){
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	clusterNamespacedName := types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}

	if finalized, err := r.ensureFinalizer(instance, clusterNamespacedName); finalized || err != nil {
		return reconcile.Result{}, err
	}

	serviceLoadBalancer, err := r.ensureLatestLoadBalancer(instance, clusterNamespacedName)

	if err != nil {
		return reconcile.Result{}, err
	}

	loadBalancerHostNames := r.extractLoadBalancerHostNames(serviceLoadBalancer)

	if len(loadBalancerHostNames) == 0 {
		return reconcile.Result{Requeue: true}, nil
	}

	err = r.ensureLatestDeployment(instance, loadBalancerHostNames, clusterNamespacedName)

	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileControlPlane) extractLoadBalancerHostNames(loadBalancer *corev1.Service)[]string{
	ips := make([]string, len(loadBalancer.Status.LoadBalancer.Ingress))

	for index,ingress := range loadBalancer.Status.LoadBalancer.Ingress{
		ips[index] = ingress.IP
	}

	/*if len(ips) == 0 {
		ips = []string{loadBalancer.Spec.ClusterIP}
	} */

	return ips
}

func (r *ReconcileControlPlane) ensureLatestDeployment(instance *caksv1alpha1.ControlPlane,
	loadBalancerHostnames []string, clusterNamespacedName types.NamespacedName)error {

	masterDeployment := &appsv1.Deployment{}

	err := r.client.Get(context.TODO(), clusterNamespacedName, masterDeployment)

	if err != nil {
		if errors.IsNotFound(err){
			return r.createMaster(clusterNamespacedName, instance, loadBalancerHostnames)
		}
		return err
	}

	desiredMasterModel, err := master.NewMaster(clusterNamespacedName, instance.Spec.ControlPlaneMaster,
		loadBalancerHostnames, master.NewResourceSplitter())

	if err != nil {
		return err
	}

	equal, err := desiredMasterModel.EqualDeployment(masterDeployment)

	if err != nil {
		return err
	}

	if !equal{
		if err = r.client.Update(context.TODO(), desiredMasterModel.BuildDeployment()); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileControlPlane) ensureLatestLoadBalancer(instance *caksv1alpha1.ControlPlane,
	clusterNamespacedName types.NamespacedName)(*corev1.Service, error){

	serviceLoadBalancer := &corev1.Service{}

	err := r.client.Get(context.TODO(), clusterNamespacedName, serviceLoadBalancer)

	if err != nil {
		if errors.IsNotFound(err){
			serviceLoadBalancer, err = r.createLoadBalancer(instance,clusterNamespacedName)
			if err != nil {
				return nil, err
			}
			return serviceLoadBalancer, nil
		}
		return nil, err
	}

	return serviceLoadBalancer, nil
}

func (r *ReconcileControlPlane) ensureFinalizer(instance *caksv1alpha1.ControlPlane,
	clusterNamespacedName types.NamespacedName)(bool,error){

	controlPlaneFinalizerName := "controlplane.finalizers.gks.globo.com"

	if instance.GetDeletionTimestamp().IsZero() {
		if !slice.ContainsString(instance.GetFinalizers(), controlPlaneFinalizerName, nil){
			instance.SetFinalizers(append(instance.GetFinalizers(),controlPlaneFinalizerName))
			if err := r.client.Update(context.TODO(), instance); err != nil {
				return false, err
			}
		}
	}else{
		if slice.ContainsString(instance.GetFinalizers(), controlPlaneFinalizerName, nil){
			masterDeployment := &appsv1.Deployment{}
			if err := r.client.Get(context.TODO(), clusterNamespacedName, masterDeployment); err == nil {
				if err := r.client.Delete(context.TODO(), masterDeployment); err != nil {
					return true, err
				}
			}else{
				if !errors.IsNotFound(err){
					return true, err
				}
			}

			instance.SetFinalizers(slice.RemoveString(instance.GetFinalizers(),controlPlaneFinalizerName, nil))
			if err := r.client.Update(context.TODO(), instance); err != nil {
				return true, err
			}
		}
		return true, nil
	}

	return false, nil
}

func (r *ReconcileControlPlane) createLoadBalancer(instance *caksv1alpha1.ControlPlane,
	namespacedName types.NamespacedName)(*corev1.Service, error){

	serviceLoadBalancer := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespacedName.Name,
			Namespace: namespacedName.Namespace,
			Labels: map[string]string{
				"app":"load-balancer",
				"cluster": namespacedName.Name,
				"tier": "control-plane",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType("LoadBalancer"),
			Ports: []corev1.ServicePort{
				{ Port: 6443, TargetPort: intstr.FromInt(6443)},
			},
			Selector: map[string]string{
				"cluster": namespacedName.Name,
			},
		},
	}

	if err := controllerutil.SetControllerReference(instance, serviceLoadBalancer, r.scheme); err != nil {
		return nil, err
	}

	if err := r.client.Create(context.TODO(), serviceLoadBalancer); err != nil {
		return nil, err
	}

	return serviceLoadBalancer, nil
}

func (r *ReconcileControlPlane) createMaster(namspacedName types.NamespacedName, instance *caksv1alpha1.ControlPlane,
	loadBalancerHostnames []string)error{
	masterModel, err := master.NewMaster(namspacedName, instance.Spec.ControlPlaneMaster, loadBalancerHostnames,
		master.NewResourceSplitter())

	if err != nil {
		return err
	}

	masterDeployment := masterModel.BuildDeployment()

	if err := controllerutil.SetControllerReference(instance, masterDeployment, r.scheme); err != nil {
		return err
	}

	if err := r.client.Create(context.TODO(), masterDeployment); err != nil {
		return err
	}

	return nil
}



