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
	edgedeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/edgedeployment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "EdgeDeployment"
	// GroupVersion is the version.
	GroupVersion string = "v1alpha1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "edgedeployments"
	// Singular is used as an alias on kubectl for display.
	Singular string = "edgedeployment"
	// CRDName is the CRD name for EdgeDeployment.
	CRDName string = Plural + "." + edgedeploy.GroupName
	// ShortName is the short alias for the CRD.
	ShortName string = "cd"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   edgedeploy.GroupName,
		Version: GroupVersion,
	}
	// SchemeBuilder is the apimachinery scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme adds to the SchemeBuilder
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&EdgeDeployment{},
		&EdgeDeploymentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
