GO_PATH=$(GOPATH)
GO_PATH?=/tmp/go
GO_SRC=$(GO_PATH)/src
PACKAGE_PATH=github.com/meritlabs
PACKAGE=$(PACKAGE_PATH)/overlord
SRC=$(GO_SRC)/$(PACKAGE)
ASSETS_OUTPUT=./pkg/controllers/statuspage/assets.go
OVERLORD_BIN=./bin/overlord
OVERSEER_BIN=./bin/overseer

GO=go
GO_FMT=$(GO) fmt

.PHONY: build
build: build-overseer build-overlord

.PHONY: build-overseer
build-overseer:
	go build -o $(OVERSEER_BIN) ./cmd/overseer/main.go

.PHONY: build-assets
build-assets:
	go-assets-builder html --package=statuspage --variable=Assets --output=$(ASSETS_OUTPUT)

.PHONY: build-overlord
build-overlord: build-assets
	go build -o $(OVERLORD_BIN) ./cmd/overlord/main.go

.PHONY: clean
clean:
	rm -rf $(OVERSEER_BIN) $(OVERLORD_BIN) $(ASSETS_OUTPUT)

.PHONY: clean-ci
clean-ci:
	rm -rf $(OVERSEER_BIN) $(OVERLORD_BIN)
	rm -rf $(SRC)

.PHONY: bootstrap
bootstrap:
	if [ ! -d "$(SRC)" ]; then mkdir -p "$(GO_SRC)/$(PACKAGE_PATH)" && ln -s "$(PWD)" "$(SRC)" ; fi
	cd "$(SRC)" && go get github.com/jessevdk/go-assets-builder
	cd "$(SRC)" && go get