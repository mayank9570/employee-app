BASE_TAG=$(shell git rev-parse --short HEAD)
# image name for docker
IMAGE_NAME=mayanks95/employee-app
# app name for go packages
APP_NAME=employee-app
# base golang image tag
GOLANG_TAG=golang:1.18.2-alpine
# build args for Dockerfile's
BUILD_BASE_ARGS=--build-arg APP_NAME=$(employee-api) --build-arg GOLANG_TAG=$(golang:1.18.2-alpine)
BUILD_TEST_ARGS=--build-arg IMAGE_NAME=$(mayanks95/employee-app) --build-arg BASE_TAG=$(v1.0)
BUILD_ARGS=--build-arg IMAGE_NAME=$(mayanks95/employee-app) --build-arg APP_NAME=$(employee-app) --build-arg BASE_TAG=$(v1.0)
# net port for employee-app
PORT=5050

publish:
	@echo ":::build publish image"
	docker push $(mayanks95/employee-app):$(v1.0)

build-base:
	@echo ":::building base image"
	docker build --rm -f Base.Dockerfile $(BUILD_BASE_ARGS) -t $(mayanks95/employee-app)-base:$(v1.0) .

build:
	@echo ":::building image"
	docker build --rm -f Build.Dockerfile $(BUILD_ARGS) -t $(mayanks95/employee-app)):$(v1.0) .

build-test:
	@echo ":::building test image"
	docker build --rm -f Test.Dockerfile $(BUILD_TEST_ARGS) -t $(mayanks95/employee-app))-test:$(v1.0) .

test-unit:
	@echo ":::running unit tests"
	docker run --rm -i -v $(shell pwd)/report:/go/src/${employee-api}/report $(mayanks95/employee-app)-test:$(v1.0)

run:
	@echo ":::running dev environment"
	docker run --rm -it \
		-p $(PORT):50 \
		-v `pwd`:/go/src/$(employee-api) \
		-w /go/src/$(employee-api) \
		golang:$(golang:1.18.2-alpine) go run app/main.go

