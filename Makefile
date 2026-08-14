GO ?= go
DOCKER ?= docker
KIND ?= kind
MINIKUBE ?= minikube
ENV ?= kind
CLUSTER_NAME ?= k8s-practice
IMAGES := hello-api gateway catalog inventory orders payments

.PHONY: test test-docker images manifests docs-check minikube-up kind-up load-images deploy smoke resilience-kind clean-kind

test:
	$(GO) test ./...

test-docker:
	$(DOCKER) run --rm -v "$(CURDIR):/src" -w /src golang:1.26.5-alpine go test ./...

images:
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/hello-api --build-arg BINARY=hello-api -t k8s-practice/hello-api:dev .
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/ecommerce/gateway --build-arg BINARY=gateway -t k8s-practice/gateway:dev .
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/ecommerce/catalog --build-arg BINARY=catalog -t k8s-practice/catalog:dev .
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/ecommerce/inventory --build-arg BINARY=inventory -t k8s-practice/inventory:dev .
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/ecommerce/orders --build-arg BINARY=orders -t k8s-practice/orders:dev .
	$(DOCKER) build -f build/Dockerfile --build-arg APP=apps/ecommerce/payments --build-arg BINARY=payments -t k8s-practice/payments:dev .

manifests:
	./tests/manifests.sh

docs-check:
	./tests/docs.sh

minikube-up:
	$(MINIKUBE) start --profile $(CLUSTER_NAME) --driver=docker --cpus=4 --memory=6144

kind-up:
	$(KIND) create cluster --name $(CLUSTER_NAME) --config clusters/kind/config.yaml

load-images:
	@if [ "$(ENV)" = "minikube" ]; then \
		for image in $(IMAGES); do $(MINIKUBE) image load --profile $(CLUSTER_NAME) k8s-practice/$$image:dev; done; \
	else \
		$(KIND) load docker-image --name $(CLUSTER_NAME) $(foreach image,$(IMAGES),k8s-practice/$(image):dev); \
	fi

deploy:
	kubectl apply -k deploy/overlays/$(ENV)
	kubectl wait --for=condition=Available deployment --all -n practice --timeout=180s
	kubectl wait --for=condition=Available deployment --all -n shop --timeout=180s

smoke:
	./tests/smoke.sh

resilience-kind:
	./tests/resilience.sh

clean-kind:
	$(KIND) delete cluster --name $(CLUSTER_NAME)
