// INTEL CONFIDENTIAL
//
// Copyright 2021-2021 Intel Corporation.
//
// This software and the related documents are Intel copyrighted materials, and your use of
// them is governed by the express license under which they were provided to you ("License").
// Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
// disclose or transmit this software or the related documents without Intel's prior written permission.
//
// This software and the related documents are provided as is, with no express or implied warranties,
// other than those that are expressly stated in the License.

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EdgeDeployment is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all
// github.com/soniabha-intc/edgedeploy/pkg/crd/edgedeployment/v1alpha1/apis
// github.com/soniabha-intc/edgedeploy/pkg/crd "edgedeployment:v1alpha1"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EdgeDeployment struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of EdgeDeployment.
	Spec EdgeDeploymentSpec `json:"spec"`
	// Observed status of EdgeDeployment.
	Status EdgeDeploymentStatus `json:"status"`
}

// EdgeDeploymentSpec is a desired state description of AppDeployment.
// +k8s:deepcopy-gen=true
type EdgeDeploymentSpec struct {
	// release_name is controller-generated name of the helm release.
	ReleaseName string `json:"releaseName"`
	// chart_uri is the URI of the helm chart used for this helm release.
	ChartURI string `json:"chartURI"`
	// namespace is the k8s namespace where the helm release should be/is installed.
	Namespace string `json:"namespace"`
	// values contains values to apply against a helm template
	ChartValues map[string]string `json:"chartValues"`
}

// EdgeDeploymentStatus describes the lifecycle status of EdgeDeployment.
type EdgeDeploymentStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// EdgeDeploymentList is the list of EdgeDeployment.
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EdgeDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of AppDeployments.
	Items []EdgeDeployment `json:"items"`
}

func (j *EdgeDeployment) String() string {
	return fmt.Sprintf(
		"\tName = %s\n\tResource Version = %s\n\tChartName = %s\n\tChartVersion = %s\n\tState = %s\n\tMessage = %s\n\t",
		j.GetName(),
		j.GetResourceVersion(),
		j.Spec.ReleaseName,
		j.Spec.ChartURI,
		j.Status.State,
		j.Status.Message,
	)
}
