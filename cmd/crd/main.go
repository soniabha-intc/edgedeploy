package main

import (
	"fmt"

	"github.com/soniabha-intc/edgedeploy/pkg/config"
	"github.com/soniabha-intc/edgedeploy/pkg/types"

	//"k8s.io/client-go/tools/clientcmd"

	appdeploy "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1"
	appdeployclient "github.com/soniabha-intc/edgedeploy/pkg/crd/appdeployment/v1alpha1/client"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	ctx := types.GetCtx()

	cfg := config.GetConfig()
	/*
		kubeconfigPath := cfg.GetKubeconfigPath()

		// Use kubeconfig to create client config.
		clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			panic(err)
		}
	*/
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Init a CRD kind.
	if _, err = appdeploy.CreateCustomResourceDefinition(apiextensionsClientSet); err != nil {
		panic(err)
	}

	// Create a CRD client interface for Jinghzhu v1.
	crdClient, err := appdeployclient.NewClient(ctx, config, cfg.GetCRDNamespace())
	if err != nil {
		panic(err)
	}

	// Create an instance of CRD.
	instanceName := "appdeployment-eii-"
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
	result, err := crdClient.CreateDefault(eiiInstance)
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else if err != nil {
		panic(err)
	}
	crdInstanceName := result.GetName()
	fmt.Println("CREATED: " + result.String())

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = crdClient.WaitForInstanceProcessed(crdInstanceName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Processed " + crdInstanceName)

	// Get the list of CRs.
	exampleList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)
}
