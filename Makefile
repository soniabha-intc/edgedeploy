export GO111MODULE = on

.PHONY: build edgedeploy
VER:=1.0

edgedeploy:
	CGO_ENABLED=0 go build -o edgedeploy ./cmd/main.go
	docker build -t edgedeploy .
install:	
	helm install deployagent helm/

uninstall:
	helm uninstall deployagent

lint:
	golangci-lint run ./pkg/ ./pkg/crd/appdeployment/ ./pkg/crd/appdeployment/v1alpha1
