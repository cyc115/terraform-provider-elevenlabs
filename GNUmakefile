VERSION   ?= 0.1.1
OS        := $(shell go env GOOS)
ARCH      := $(shell go env GOARCH)
INSTALL_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/cyc115/elevenlabs/$(VERSION)/$(OS)_$(ARCH)
BINARY    := terraform-provider-elevenlabs

.PHONY: install test lint docs docs-check

install:
	go build -o $(BINARY) .
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

test:
	go test ./... -v

lint:
	go vet ./...

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name elevenlabs

docs-check:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name elevenlabs
	git diff --exit-code docs/ || (echo "docs/ out of sync — run 'make docs' and commit"; exit 1)
