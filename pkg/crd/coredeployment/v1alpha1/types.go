package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CoreDeployment is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1/apis github.com/soniabha-intc/edgedeploy/pkg/crd "appdeployment:v1alpha1"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=appdeployment
type CoreDeployment struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of AppDeployment.
	Spec CoreDeploymentSpec `json:"spec"`
	// Observed status of AppDeployment.
	Status CoreDeploymentStatus `json:"status"`
}

// CoreDeploymentSpec is a desired state description of AppDeployment.
// +k8s:deepcopy-gen=true
type CoreDeploymentSpec struct {
	// ChartName is the name of the helm chart to deploy.
	ChartName string `json:"chartName"`
	// ChartVersion is the version of the helm chart to deploy.
	ChartVersion string `json:"chartVersion"`
	// Parameters is key/value pairs of helm chart values.yaml
	Parameters map[string]string `json:"parameters"`
}

// CoreDeploymentStatus describes the lifecycle status of AppDeployment.
type CoreDeploymentStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// CoreDeploymentList is the list of AppDeployment.
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=appdeployment
type CoreDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of AppDeployments.
	Items []CoreDeployment `json:"items"`
}

func (j *CoreDeployment) String() string {
	return fmt.Sprintf(
		"\tName = %s\n\tResource Version = %s\n\tChartName = %s\n\tChartVersion = %s\n\tState = %s\n\tMessage = %s\n\t",
		j.GetName(),
		j.GetResourceVersion(),
		j.Spec.ChartName,
		j.Spec.ChartVersion,
		j.Status.State,
		j.Status.Message,
	)
}
