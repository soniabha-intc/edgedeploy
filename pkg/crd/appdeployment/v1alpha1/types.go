package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppDeployment is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/soniabha-intc/edgedeploy/KubernetesCRD/pkg/crd/appdeployment/v1alpha1/apis github.com/soniabha-intc/edgedeploy/KubernetesCRD/pkg/crd "appdeployment:v1alpha1"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=appdeployment
type AppDeployment struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of AppDeployment.
	Spec AppDeploymentSpec `json:"spec"`
	// Observed status of AppDeployment.
	Status AppDeploymentStatus `json:"status"`
}

// AppDeploymentStatus is a desired state description of AppDeployment.
// +k8s:deepcopy-gen=true
type AppDeploymentSpec struct {
	// ChartName is the name of the helm chart to deploy.
	ChartName string `json:"chartName"`
	// ChartVersion is the version of the helm chart to deploy.
	ChartVersion string `json:"chartVersion"`
	// PodList is key/value pairs of helm chart values.yaml
	Parameters map[string]string `json:"parameters"`
	// fqdn is optional for app, If present A record is added into EdgeDNS
	Fqdns map[string]string `json:"fqdns"`
}

// AppDeploymentStatus describes the lifecycle status of AppDeployment.
type AppDeploymentStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// AppDeploymentList is the list of AppDeployment.
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu
type AppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of Jinghzhus.
	Items []AppDeployment `json:"items"`
}

func (j *AppDeployment) String() string {
	return fmt.Sprintf(
		"\tName = %s\n\tResource Version = %s\n\tChartName = %s\n\tChartVersion = %s\n\tState = %s\n\tMessage = %s\n\t",
		j.GetName(),
		j.GetResourceVersion(),
		j.Spec.ChartName,
		j.Spec.ChartVersion,
		//strings.Join(j.Spec.PodList, ", "),
		j.Status.State,
		j.Status.Message,
	)
}
