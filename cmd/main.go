package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	deploy "github.com/soniabha-intc/edgedeploy/pkg"
	"github.com/soniabha-intc/edgedeploy/pkg/config"
	"k8s.io/client-go/rest"
)

const cfgPath = "configs/edgedeploy.json"

func main() {

	parenCtx, cancel := context.WithCancel(context.Background())
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		fmt.Printf("Received signal: %#v", sig)
		cancel()
	}()

	err := config.InitConfig(cfgPath)
	if err != nil {
		fmt.Printf("InitConfig failed %s", err.Error())
		return
	}
	// creates the in-cluster config
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("InClusterConfig failed %s", err.Error())
		return
	}

	err = deploy.AppCRD(parenCtx, kubeconfig)
	if err != nil {
		fmt.Println("CRD AppDeployment failed")
		return
	}

	err = deploy.CoreCRD(parenCtx, kubeconfig)
	if err != nil {
		fmt.Println("CRD CoreDeployment failed")
		return
	}

	//ticker := time.NewTicker(5 * time.Second)
	deploy.StartReconcileLoop(parenCtx)
	//done := make(chan bool)
	//go deploy.Reconcile(parenCtx, done, ticker)
	time.Sleep(9 * time.Second)
	deploy.DeployReconcile.Done <- true
}
