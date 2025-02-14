# Variables
DOCKER_USERNAME = seanminingah
IMAGE_NAME = sil-backend-assessment
TAG ?= latest

.PHONY: generate clean

build:
	@echo "Building Docker image with tag: $(TAG)..."
	@docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG) .

tag:
	@echo "Tagging Docker image with tag: $(TAG)..."
	docker tag $(IMAGE_NAME):$(TAG) $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)

push:
	@echo "Pushing Docker image with tag: $(TAG)..."
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)

clean:
	@echo "Cleaning Docker image with tag: $(TAG)..."
	docker rmi $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)

generate:
	@echo "Generating GraphQL code..."
	@go generate ./...
	# @go run github.com/99designs/gqlgen generate

clean:
	@echo "Cleaning generated GraphQL files..."
	@rm -rf internal/adapters/handlers/graphql/generated

test:
	@echo "Running tests..."
	@go test ./...

start:
	@echo "Starting server..."
	@go run cmd/main.go
