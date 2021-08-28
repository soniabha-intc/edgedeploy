package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/soniabha-intc/edgedeploy/pkg/config"
)

type EdgeReconcile struct {
	Ticker *time.Ticker
	Done   chan bool
}

var DeployReconcile EdgeReconcile

func StartReconcileLoop(ctx context.Context) {
	duration := config.GetConfig().ReconcileTimer
	DeployReconcile.Ticker = time.NewTicker(time.Duration(duration) * time.Second)
	DeployReconcile.Done = make(chan bool)
	go Reconcile(ctx, DeployReconcile.Done, DeployReconcile.Ticker)

}

func StopReconcileLoop() {
	DeployReconcile.Ticker.Stop()
	DeployReconcile.Done <- true
}

func Reconcile(ctx context.Context, done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
			GRPCCoreReplacement(ctx)
			//
		}
	}
}
