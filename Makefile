export GO111MODULE = on

.PHONY: build edgedeploy
VER:=1.0

edgedeploy:
	CGO_ENABLED=0 go build -o edgedeploy ./cmd/crd/main.go
	#docker build -t edgedeploy .
	#helm install edgedeploy helm/
