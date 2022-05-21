HOSTNAME = kentik
NAMESPACE = automation
NAME = kentik-cloudexport
BINARY = terraform-provider-${NAME}
VERSION := $(shell echo `git tag --list 'v*' | tail -1 | cut -d v -f 2` | sed -e 's/^$$/0.1.0/')
OS_ARCH := $(shell printf "%s_%s" `go env GOHOSTOS` `go env GOHOSTARCH`)

default: install

build:
	go mod tidy
	go build -o ${BINARY}

check-docs:
	./tools/check_docs.sh

check-go-mod:
	./tools/check_go_mod.sh

docs:
	go generate

fmt:
	./tools/fmt.sh

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

lint:
	golangci-lint run

test:
	go test ./... -timeout=5m

.PHONY: build check-docs docs fmt install lint test
