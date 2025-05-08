PROJECT := light
# Where to push the docker image.
REGISTRY ?= yazhivotnoe
IMAGE := $(REGISTRY)/$(PROJECT)
VERSION ?= 0.0.1

.PHONY: build deps test image push rm_image

build: 
	rm -f ./lalachka && go build -o lalachka ./cmd/main.go

run: 
	./lalachka

deps:
	go mod download && go mod tidy && go mod verify

test:
	go test -v ./...

list_pkgs:
	go list ./cmd/... ./internal/...

image:
	docker build --platform=linux/amd64 -t $(IMAGE):$(VERSION) . -f Dockerfile

rm_image: 

push:
	docker push $(IMAGE):$(VERSION)

	docker rmi $(IMAGE):$(VERSION) -f
