package client

import (
	"encoding/json"
	"fmt"
	"time"

	coredeployment "github.com/soniabha-intc/edgedeploy/pkg/crd/coredeployment/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	"github.com/soniabha-intc/edgedeploy/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForInstanceProcessed is used for monitor the creation of a CRD instance.
func (c *Client) WaitForInstanceProcessed(name string) error {
	return wait.Poll(time.Second, 3*time.Second, func() (bool, error) {
		instance, err := c.Get(name, metav1.GetOptions{})
		if err == nil && instance.Status.State == types.StatePending {
			return true, nil
		}
		fmt.Printf("Fail to wait for CRD instance processed: %+v\n", err)

		return false, err
	})
}

// Create post an instance of CRD into Kubernetes with given create options.
func (c *Client) Create(obj *coredeployment.CoreDeployment, opts metav1.CreateOptions) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Create(c.GetContext(), obj, opts)
}

// CreateDefault post an instance of CRD into Kubernetes without create options.
func (c *Client) CreateDefault(obj *coredeployment.CoreDeployment) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Create(c.GetContext(), obj, metav1.CreateOptions{})
}

// Update puts new instance of CRD to replace the old one by given update options.
func (c *Client) Update(obj *coredeployment.CoreDeployment, opts metav1.UpdateOptions) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Update(c.GetContext(), obj, opts)
}

// UpdateDefault puts new instance of CRD to replace the old one without update options.
func (c *Client) UpdateDefault(obj *coredeployment.CoreDeployment) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Update(c.GetContext(), obj, metav1.UpdateOptions{})
}

// UpdateSpecAndStatus updates the spec and status filed of CRD.
// If only want to update some sub-resource, please use Patch instead.
func (c *Client) UpdateSpecAndStatus(name string, appdeployspec *coredeployment.CoreDeploymentSpec, appdeploystatus *coredeployment.CoreDeploymentStatus) (*coredeployment.CoreDeployment, error) {
	instance, err := c.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	instance.Spec = *appdeployspec
	instance.Status = *appdeploystatus

	return c.Update(instance, metav1.UpdateOptions{})
}

// Patch applies the patch and returns the patched  instance.
func (c *Client) Patch(name string, pt apimachinerytypes.PatchType, data []byte, subresources ...string) (*coredeployment.CoreDeployment, error) {
	var result coredeployment.CoreDeployment
	err := c.clientset.RESTClient().Patch(pt).
		Namespace(c.namespace).
		Resource(c.plural).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(c.GetContext()).
		Into(&result)

	return &result, err
}

// PatchJSONType uses JSON Type (RFC6902) in PATCH.
func (c *Client) PatchJSONType(name string, ops []PatchJSONTypeOps) (*coredeployment.CoreDeployment, error) {
	patchBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Patch(c.GetContext(), name, apimachinerytypes.JSONPatchType, patchBytes, metav1.PatchOptions{})
}

// PatchSpec only updates the spec field of AppDeployment v1alpha1, which is /spec.
func (c *Client) PatchSpec(name string, appdeployspec *coredeployment.CoreDeploymentSpec) (*coredeployment.CoreDeployment, error) {
	ops := make([]PatchJSONTypeOps, 1)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/spec"
	ops[0].Value = appdeployspec

	return c.PatchJSONType(name, ops)
}

// PatchStatus only updates the status field of AppDeployment v1alpha1, which is /status.
func (c *Client) PatchStatus(name string, appdeploystatus *coredeployment.CoreDeploymentStatus) (*coredeployment.CoreDeployment, error) {
	ops := make([]PatchJSONTypeOps, 1)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/status"
	ops[0].Value = appdeploystatus

	return c.PatchJSONType(name, ops)
}

// PatchSpecAndStatus performs patch for both spec and status field of AppDeployment.
func (c *Client) PatchSpecAndStatus(name string, appdeployspec *coredeployment.CoreDeploymentSpec, appdeploystatus *coredeployment.CoreDeploymentStatus) (*coredeployment.CoreDeployment, error) {
	ops := make([]PatchJSONTypeOps, 2)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/spec"
	ops[0].Value = appdeployspec
	ops[1].Op = PatchJSONTypeReplace
	ops[1].Path = "/status"
	ops[1].Value = appdeploystatus

	return c.PatchJSONType(name, ops)
}

// Delete removes the CRD instance by given name and delete options.
func (c *Client) Delete(name string, opts metav1.DeleteOptions) error {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Delete(c.GetContext(), name, opts)
}

// DeleteDefault removes the CRD instance without delete options.
func (c *Client) DeleteDefault(name string) error {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Delete(c.GetContext(), name, metav1.DeleteOptions{})
}

// Get returns a pointer to the CRD instance.
func (c *Client) Get(name string, opts metav1.GetOptions) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Get(c.GetContext(), name, opts)
}

// GetDefault retrieves the crd instance without get options.
func (c *Client) GetDefault(name string) (*coredeployment.CoreDeployment, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).Get(c.GetContext(), name, metav1.GetOptions{})
}

// List returns a list of CRD instances by given list options.
func (c *Client) List(opts metav1.ListOptions) (*coredeployment.CoreDeploymentList, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).List(c.GetContext(), opts)
}

// ListDefaultDefault returns a list of CRD instances without list options.
func (c *Client) ListDefaultDefault() (*coredeployment.CoreDeploymentList, error) {
	return c.clientset.CoredeploymentV1alpha1().CoreDeployments(c.namespace).List(c.GetContext(), metav1.ListOptions{})
}
