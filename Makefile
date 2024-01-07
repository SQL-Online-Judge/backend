# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint
BINARY_NAME=core

all: lint build test
clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
lint:
	$(GOLINT) run ./...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
