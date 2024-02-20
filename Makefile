# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint
PROJECT_NAME=SQL-Online-Judge
SERVICES=core judger
OUT_DIR=./bin

all: clean lint $(addprefix $(OUT_DIR)/, $(SERVICES)) test
clean:
	$(GOCMD) clean
	rm -rf $(OUT_DIR)
lint:
	$(GOLINT) run ./...
$(OUT_DIR)/%: clean
	$(GOBUILD) -o $(OUT_DIR)/$* -v ./cmd/$*
test:
	$(GOTEST) -v ./...
debug:
	docker compose down
	docker compose build
	docker compose up
