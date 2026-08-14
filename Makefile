GO ?= go
DOCKER ?= docker

.PHONY: test test-docker images manifests docs-check

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

