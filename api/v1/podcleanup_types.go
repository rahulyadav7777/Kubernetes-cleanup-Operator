package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// PodCleanupSpec defines the desired state of PodCleanup
type PodCleanupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// PodCleanupStatus defines the observed state of PodCleanup
type PodCleanupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodCleanup is the Schema for the podcleanups API
type PodCleanup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCleanupSpec   `json:"spec,omitempty"`
	Status PodCleanupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodCleanupList contains a list of PodCleanup
type PodCleanupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCleanup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCleanup{}, &PodCleanupList{})
}
