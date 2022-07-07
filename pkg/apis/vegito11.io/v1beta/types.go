package v1beta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true
// +kubebuilder:resource:shortName=awth

// Our Custom Object structure
type AWSAuthMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthSpec `json:"spec"`
}

type AuthSpec struct {
	MapRoles []MapRole `json:"rolesmap,omitempty"`
	MapUsers []MapUser `json:"usersmap,omitempty"`
}

type MapRole struct {
	RoleARN  string   `json:"rolearn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

type MapUser struct {
	UserARN  string   `json:"userarn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:deepcopy-gen=true

// List of our CR
type AWSAuthMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AWSAuthMap `json:"items"`
}
