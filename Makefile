# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: deps test

test:
		$(GOTEST) -v -cover ./...

deps:
		$(GOGET) -v ./...
    