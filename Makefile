run:
	go run cmd/main.go

# Define variables
IMAGE_NAME := backend-report-service
DOCKER_REPO := stanthikun802
DOCKER_TAG := latest

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Tag the Docker image
tag:
	docker tag $(IMAGE_NAME) $(DOCKER_REPO)/$(IMAGE_NAME):$(DOCKER_TAG)

# Push the Docker image to the repository
push:
	docker push $(DOCKER_REPO)/$(IMAGE_NAME):$(DOCKER_TAG)

# Remove the locally built Docker image
clean:
	docker rmi $(IMAGE_NAME)

# Chain tasks together
all: build tag push clean

# Define phony targets to avoid conflicts with file names
.PHONY: build tag push clean all