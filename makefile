SHELL := /bin/bash

# curl --user "admin@example.com:gophers" http://localhost:3000/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
# export TOKEN= TOKEN BODY
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/5cf37266-3473-4006-984f-9325122678b7
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/1/1

all: core metrics

core:
	docker build \
	-f zarf/docker/dockerfile.shop \
	-t shop-amd64:1.0 \
	--build-arg VSF_REF=`git rev-parse HEAD` \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.


metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t metrics-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.


# ==============================================================================
# Running from within k8s/dev
kind-up:
	kind create cluster --image kindest/node:v1.19.4 --name starter-cluster --config zarf/k8s/dev/kind-config.yaml


kind-down:
	kind delete cluster --name starter-cluster

kind-load:
	kind load docker-image shop-amd64:1.0 --name starter-cluster
	# kind load docker-image metrics-amd64:1.0 --name starter-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -


kind-status:
	kubectl get nodes
	kubectl get pods --watch


kind-status-full:
	kubectl describe pod -lapp=shop
kind-update: core
	kind load docker-image shop-amd64:1.0 --name starter-cluster
	kubectl delete pods -lapp=shop


kind-logs:
	kubectl logs -lapp=shop --all-containers=true -f
# ==============================================================================

run:
	go run ./cmd/app/main.go


runa:
	go run ./cmd/admin/main.go


lint:
	export CGO_ENABLED=0 && go test -v ./...