VERSION   ?= 0.1.0
OS        := $(shell go env GOOS)
ARCH      := $(shell go env GOARCH)
INSTALL_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/cyc115/elevenlabs/$(VERSION)/$(OS)_$(ARCH)
BINARY    := terraform-provider-elevenlabs

.PHONY: install test lint

install:
	go build -o $(BINARY) .
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

test:
	go test ./... -v

lint:
	go vet ./...
