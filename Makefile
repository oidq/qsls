# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=qsls
VERSION := $(shell git describe --tags)

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.version=${VERSION}" -v
test:
		$(GOTEST)  ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
