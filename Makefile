GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -count=1
GOGET=$(GOCMD) get
BINARY_NAME=go-currency-exchange

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

deps:
	$(GOGET) -u ./...

dev:
	air

docs:
	swag init -g cmd/server/main.go

.PHONY: all build clean test run deps dev docs