# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOINSTALL=$(GOCMD) install
BINARY_NAME=grpc-crud

all: deps build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
proto:
	protoc --go_out="plugins=grpc:." hirepb/hire.proto
install:
	$(GOINSTALL)
run:
	$(GORUN) ./*.go
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
deps:
	$(GOGET) github.com/mongodb/mongo-go-driver/bson
	$(GOGET) github.com/mongodb/mongo-go-driver/bson/primitive	
	$(GOGET) github.com/mongodb/mongo-go-driver/mongo
	$(GOGET) google.golang.org/grpc/codes
	$(GOGET) google.golang.org/grpc/status
	$(GOGET) google.golang.org/grpc
	$(GOGET) google.golang.org/grpc/reflection