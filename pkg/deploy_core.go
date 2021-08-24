package deploy

import (
	"context"
	"fmt"

	"github.com/soniabha-intc/edgedeploy/pkg/config"
	"github.com/soniabha-intc/edgedeploy/pkg/types"

	coredeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/coredeployment/v1alpha"
	coredeployclient "github.com/soniabha-intc/edgedeploy/pkg/crd/coredeployment/v1alpha/client"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var clientCore *coredeployclient.Client

// CoresDeployed is the structure for storing the coredeployment CRDs
type CoresDeployed struct {
	coresDeployed map[string]*coredeploy.CoreDeployment
}

// CoresDeployed stores the Apps which are deployed in the cluster and the count
var CoreDeployed CoresDeployed

// AppCRD invokes the AppDeployment CRD creation and CR CRUD
func CoreCRD(parentCtx context.Context, kubeconfig *rest.Config) error {

	//ctx := types.GetCtx()

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(kubeconfig)
	if err != nil {
		//panic(err)
		return err
	}

	// Create the AppDeployment kind CRD
	if _, err = coredeploy.CreateCustomResourceDefinitionCoreDeploy(apiextensionsClientSet); err != nil {
		fmt.Println("CreateCustomResourceDefinitionCoreDeploy failed")
		return err
	}

	// Create the AppDeployment Client set
	err = InitCoreDeployClient(parentCtx, kubeconfig)
	if err != nil {
		fmt.Println("Failure in InitCoreDeployClient")
	}

	return err

}

func GRPCCoreReplacement(parentCtx context.Context) error {
	// TODO the following APIs to be invoked on gRPC responses
	// Create the App Deployment CR
	instanceName := "coredeployment-druid-"
	err := CreateCoreDeploymentCR(parentCtx, instanceName)
	if err != nil {
		fmt.Printf("Creating CR %s failed", instanceName)
		return err
	}
	// Update  the App Deployment CR

	for _, value := range CoreDeployed.coresDeployed {

		value.Spec.ChartVersion = "2.0.0"

		err = UpdateCoreDeploymentCR(value)
		if err != nil {
			fmt.Printf("Updating CR %s failed", value.ObjectMeta.GenerateName)
			return err
		}

	}
	// Delete the App Deployment CR

	for key := range CoreDeployed.coresDeployed {

		err = DeleteCoreDeploymentCR(key)
		if err != nil {
			fmt.Printf("Deleting CR %s failed", key)
			return err
		}
	}
	return nil
}

// InitAppDeployClient creates the client for AppDeployment CRD
func InitCoreDeployClient(ctx context.Context, kubeconfig *rest.Config) error {

	cfg := config.GetConfig()
	// Create a CRD client interface for CoreDeployment v1alpha1.
	crdClient, err := coredeployclient.NewClient(ctx, kubeconfig, cfg.AgentNamespace)
	if err != nil {
		//panic(err)
		return err
	}

	clientCore = crdClient
	return nil
}

// CreateAppDeploymentCR creates the AppDeployment CR instance
func CreateCoreDeploymentCR(ctx context.Context, instanceName string) error {

	// TODO Get the App Deployment Params from Controller
	// Create an instance of CRD.

	druidInstance := &coredeploy.CoreDeployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: instanceName,
		},
		Spec: coredeploy.CoreDeploymentSpec{
			ChartName:    "druid",
			ChartVersion: "1.0.0",
			Parameters:   make(map[string]string),
		},
		Status: coredeploy.CoreDeploymentStatus{
			State:   types.StatePending,
			Message: "Created but not processed yet",
		},
	}
	druidInstance.Spec.Parameters["sriov.enableSriov"] = "true"

	result, err := clientCore.CreateDefault(druidInstance)
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else if err != nil {
		panic(err)
	}

	fmt.Println("CREATED: " + result.String())
	if CoreDeployed.coresDeployed == nil {
		CoreDeployed.coresDeployed = make(map[string]*coredeploy.CoreDeployment, 10)
	}
	crdInstanceName := result.GetName()
	CoreDeployed.coresDeployed[crdInstanceName] = result
	return nil
}

// DeleteAppDeploymentCR deletes the AppDeployment CR
func DeleteCoreDeploymentCR(crdInstanceName string) error {

	err := clientCore.DeleteDefault(crdInstanceName)
	delete(CoreDeployed.coresDeployed, crdInstanceName)
	fmt.Println("DELETED: " + crdInstanceName)

	return err
}

// UpdateAppDeploymentCR updates the AppDeployment CR
func UpdateCoreDeploymentCR(obj *coredeploy.CoreDeployment) error {

	result, err := clientCore.UpdateDefault(obj)
	if err != nil {
		fmt.Printf("UpdateCoreDeployCRD failed %s", result.String())
	}
	fmt.Println("UPDATED: " + result.String())
	return nil

}
