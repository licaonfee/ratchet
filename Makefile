# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: deps utest grep

deps:
		$(GOGET) github.com/DATA-DOG/go-sqlmock
		$(GOGET) -v ./...

utest:
		$(GOTEST) -v -cover ./...
	
grep:
		cd examples/grep ; \
		echo "barcelona" | $(GOCMD) run main.go "^bar" 