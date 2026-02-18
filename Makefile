BINARY=terraform-provider-manta
VERSION=0.0.1
HOSTNAME=registry.terraform.io
NAMESPACE=gagno
NAME=manta

OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

ifeq ($(OS),windows)
  EXT=.exe
  PLUGIN_DIR=$(subst \,/,$(APPDATA))/terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS)_$(ARCH)
else
  EXT=
  PLUGIN_DIR=$(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS)_$(ARCH)
endif

# Binary name required by filesystem_mirror unpacked layout
INSTALL_BINARY=$(BINARY)_v$(VERSION)$(EXT)

default: fmt lint install generate

build:
	go build -o $(BINARY)$(EXT)

ifeq ($(OS),windows)
install: build
	powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(PLUGIN_DIR)' | Out-Null"
	powershell -NoProfile -Command "Copy-Item -Force '$(BINARY)$(EXT)' '$(PLUGIN_DIR)/$(INSTALL_BINARY)'"
else
install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY)$(EXT) $(PLUGIN_DIR)/$(INSTALL_BINARY)
endif

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: default build install fmt lint generate test testacc
