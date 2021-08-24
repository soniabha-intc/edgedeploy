package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/soniabha-intc/edgedeploy/pkg/config"
	appdeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1"
	appdeploymentclientset "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1/apis/clientset/versioned"
	"k8s.io/client-go/rest"

	"github.com/soniabha-intc/edgedeploy/pkg/types"
)

const (
	// PatchJSONTypeReplace is the command to replace
	PatchJSONTypeReplace string = "replace"
	// PatchJSONTypeAdd is the command to add
	PatchJSONTypeAdd string = "add"
)

var (
	onceDefaultAppdeployV1Alpha1Client sync.Once
	defaultClient                      *Client
	//validPatchResources                map[string]string
)

// Client is an API client to help perform CRUD for CRD instances.
type Client struct {
	clientset *appdeploymentclientset.Clientset
	namespace string
	plural    string
	ctx       context.Context
}

// PatchJSONTypeOps describes the operations for PATCH defined in RFC6902. https://tools.ietf.org/html/rfc6902
// The supported operations are: add, remove, replace, move, copy and test.
// When we news a AppDeployment instance, we'll set default value for all fields. So, when you want to patch a AppDeployment,
// DO NOT use remove. Please use replace, even if you want to keep that field "empty".
// Example:
// 	things := make([]IntThingSpec, 2)
// 	things[0].Op = "replace"
// 	things[0].Path = "/status/message"
// 	things[0].Value = "1234"
// 	things[1].Op = "replace"
// 	things[1].Path = "/status/state"
// 	things[1].Value = ""
type PatchJSONTypeOps struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// GetNamespace returns the namespace the client talks to.
func (c *Client) GetNamespace() string {
	return c.namespace
}

// GetPlural returns the plural the client is managing.
func (c *Client) GetPlural() string {
	return c.plural
}

// GetContext returns the context of client.
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// CreateAppDeploymentClientset returns the clientset for CRD AppDeployment v1alpha1 in singleton way.
func CreateAppDeploymentClientset(config *rest.Config) (*appdeploymentclientset.Clientset, error) {
	/*restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}*/
	clientset, err := appdeploymentclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// NewClient accepts kubeconfig path and namespace. Return the API client interface for CRD AppDeployment v1alpha1.
func NewClient(ctx context.Context, config *rest.Config, namespace string) (*Client, error) {
	clientset, err := CreateAppDeploymentClientset(config)
	if err != nil {
		fmt.Printf("Fail to init CRD API clientset for Jinghuazhu v1: %+v\n", err.Error())

		return nil, err
	}
	c := &Client{
		clientset: clientset,
		namespace: namespace,
		plural:    appdeploy.Plural,
		ctx:       ctx,
	}

	return c, nil
}

// GetDefaultClient returns an API client interface for CRD AppDeployment v1alpha1. It assumes the kubeconfig
// is available at default path and the target CRD namespace is the default namespace.
func GetDefaultClient() *Client {
	onceDefaultAppdeployV1Alpha1Client.Do(func() {

		kconfig, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		cfg := config.GetConfig()
		clientset, err := CreateAppDeploymentClientset(kconfig)
		if err != nil {
			panic("Fail to init default CRD API client for AppDeployment v1alpha1: " + err.Error())
		}
		defaultClient = &Client{
			clientset: clientset,
			namespace: cfg.AgentNamespace,
			plural:    appdeploy.Plural,
			ctx:       types.GetCtx(),
		}
	})

	return defaultClient
}
