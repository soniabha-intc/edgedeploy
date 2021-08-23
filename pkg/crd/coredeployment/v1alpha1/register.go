package v1alpha1

import (
	coredeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/coredeployment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "CoreDeployment"
	// GroupVersion is the version.
	GroupVersion string = "v1alpha1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "coredeployments"
	// Singular is used as an alias on kubectl for display.
	Singular string = "coredeployment"
	// CRDName is the CRD name for CoreDeployment.
	CRDName string = Plural + "." + coredeploy.GroupName
	// ShortName is the short alias for the CRD.
	ShortName string = "ad"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   coredeploy.GroupName,
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
		&CoreDeployment{},
		&CoreDeploymentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
