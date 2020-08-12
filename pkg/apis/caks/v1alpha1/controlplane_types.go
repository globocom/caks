package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControlPlaneMaster struct {
	Count int `json:"count,omitempty"`
	ControlPlaneMasterSettings `json:"settings,omitempty"`
	corev1.ResourceRequirements `json:"resources,omitempty"`
}

type ControlPlaneMasterSettings struct {
	MasterSecretName string `json:"certsSecret,omitempty"`
	AdmissionPlugins []string `json:"admissionPlugins,omitempty"`
	ServiceClusterIPRange string `json:"serviceClusterIpRange,omitempty"`
	ClusterCIDR string `json:"clusterCidr,omitempty"`
	EncryptionSecretName string `json:"encryptionSecret,omitempty"`
}

// ControlPlaneSpec defines the desired state of ControlPlane
type ControlPlaneSpec struct {
	ControlPlaneMaster `json:"master,omitempty"`
}

// ControlPlaneStatus defines the observed state of ControlPlane
type ControlPlaneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControlPlane is the Schema for the controlplanes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=controlplanes,scope=Namespaced
type ControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlPlaneSpec   `json:"spec,omitempty"`
	Status ControlPlaneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControlPlaneList contains a list of ControlPlane
type ControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ControlPlane{}, &ControlPlaneList{})
}
