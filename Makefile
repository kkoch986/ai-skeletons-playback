SHELL=/bin/sh
PACKAGES=$(shell go list ./...)
BIN_DIR=bin
BIN_NAME=playback
PACKAGE_NAME=$(shell go list .)
SOURCE_DIR=.
SOURCES=$(shell find $(SOURCE_DIR) -name '*.go')

GIT_SHA=$(shell git rev-parse --short=8 HEAD)

all: docker-build

build: fmt $(BIN_DIR)/$(BIN_NAME)

check-fmt:
	@gofmt -d . | read; if [ $$? == 0 ]; then echo "gofmt check failed for:"; gofmt -d -l .; exit 1; fi

fmt:
	go fmt ./...

$(BIN_DIR)/$(BIN_NAME): $(SOURCES)
	mkdir -p $(BIN_DIR)
	go build -x -o $(BIN_DIR)/$(BIN_NAME) 

test:
	go test -count=1 -v $(PACKAGES)

clean:
	rm -rf $(BIN_DIR)
	rm -f .coverage*

go-lint:
	$(eval GOLINT_INSTALLED := $(shell which golangci-lint))
	@if [ "$(GOLINT_INSTALLED)" = "" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1; \
	fi;

lint: go-lint
	golangci-lint run

docker-build: build
	docker build . \
		-t kkoch986/ai-skeletons-playback:latest \
		#-t kkoch986/ai-skeletons-playback:$(GIT_SHA);

.PHONY: clean lint all build fmt test go-lint docker-build
