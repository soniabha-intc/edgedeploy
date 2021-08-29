package deploy

import (
	"context"
	"fmt"

	"github.com/soniabha-intc/edgedeploy/pkg/config"
	"github.com/soniabha-intc/edgedeploy/pkg/types"

	appdeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1"
	appdeployclient "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1/client"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var client *appdeployclient.Client

// AppsDeployed is the structure for storing the appdeployment CRDs
type AppsDeployed struct {
	appsDeployed map[string]*appdeploy.AppDeployment
}

// ApplicationsDeployed stores the Apps which are deployed in the cluster and the count
var ApplicationsDeployed AppsDeployed

// AppCRD invokes the AppDeployment CRD creation and CR CRUD
func AppCRD(parentCtx context.Context, kubeconfig *rest.Config) error {

	//ctx := types.GetCtx()

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(kubeconfig)
	if err != nil {
		//panic(err)
		return err
	}

	// Create the AppDeployment kind CRD
	if _, err = appdeploy.CreateCustomResourceDefinitionAppDeploy(apiextensionsClientSet); err != nil {
		fmt.Println("CreateCustomResourceDefinitionAppDeploy failed")
		return err
	}

	// Create the AppDeployment Client set
	err = InitAppDeployClient(parentCtx, kubeconfig)
	if err != nil {
		fmt.Println("Failure in InitAppDeployClient")
	}

	return err

}

func GRPCReplacement(parentCtx context.Context) error {
	// TODO the following APIs to be invoked on gRPC responses
	// Create the App Deployment CR
	instanceName := "appdeployment-eii-"
	err := CreateAppDeploymentCR(parentCtx, instanceName)
	if err != nil {
		fmt.Printf("Creating CR %s failed", instanceName)
		return err
	}
	// Update  the App Deployment CR

	for _, value := range ApplicationsDeployed.appsDeployed {

		value.Spec.ChartVersion = "2.0.0"

		err = UpdateAppDeploymentCR(value)
		if err != nil {
			fmt.Printf("Updating CR %s failed", value.ObjectMeta.GenerateName)
			return err
		}

	}
	// Delete the App Deployment CR
	/*
		for key := range ApplicationsDeployed.appsDeployed {

			err = DeleteAppDeploymentCR(key)
			if err != nil {
				fmt.Printf("Deleting CR %s failed", key)
				return err
			}
		}*/
	return nil
}

// InitAppDeployClient creates the client for AppDeployment CRD
func InitAppDeployClient(ctx context.Context, kubeconfig *rest.Config) error {

	cfg := config.GetConfig()
	// Create a CRD client interface for AppDeployment v1alpha1.
	crdClient, err := appdeployclient.NewClient(ctx, kubeconfig, cfg.AgentNamespace)
	if err != nil {
		//panic(err)
		return err
	}

	client = crdClient
	return nil
}

// CreateAppDeploymentCR creates the AppDeployment CR instance
func CreateAppDeploymentCR(ctx context.Context, instanceName string) error {

	// TODO Get the App Deployment Params from Controller
	// Create an instance of CRD.

	eiiInstance := &appdeploy.AppDeployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: instanceName,
		},
		Spec: appdeploy.AppDeploymentSpec{
			ChartName:    "eii-app",
			ChartVersion: "1.0.0",
			Parameters:   make(map[string]string),
			Fqdns:        make(map[string]string),
		},
		Status: appdeploy.AppDeploymentStatus{
			State:   types.StatePending,
			Message: "Created but not processed yet",
		},
	}
	eiiInstance.Spec.Parameters["namespace"] = "app_namespace"
	eiiInstance.Spec.Fqdns["name"] = "app_fqdn"

	result, err := client.CreateDefault(eiiInstance)
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else if err != nil {
		panic(err)
	}

	fmt.Println("CREATED: " + result.String())
	if ApplicationsDeployed.appsDeployed == nil {
		ApplicationsDeployed.appsDeployed = make(map[string]*appdeploy.AppDeployment, 10)
	}
	crdInstanceName := result.GetName()
	ApplicationsDeployed.appsDeployed[crdInstanceName] = result
	return nil
}

// DeleteAppDeploymentCR deletes the AppDeployment CR
func DeleteAppDeploymentCR(crdInstanceName string) error {

	err := client.DeleteDefault(crdInstanceName)
	delete(ApplicationsDeployed.appsDeployed, crdInstanceName)
	fmt.Println("DELETED: " + crdInstanceName)

	return err
}

// UpdateAppDeploymentCR updates the AppDeployment CR
func UpdateAppDeploymentCR(obj *appdeploy.AppDeployment) error {

	result, err := client.UpdateDefault(obj)
	if err != nil {
		fmt.Printf("UpdateAppDeployCRD failed %s", result.String())
	}
	fmt.Println("UPDATED: " + result.String())
	return nil

}
