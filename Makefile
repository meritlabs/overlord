GO_PATH=$(GOPATH)
GO_PATH?=/tmp/go
GO_SRC=$(GO_PATH)/src
PACKAGE_PATH=github.com/meritlabs
PACKAGE=$(PACKAGE_PATH)/overlord
SRC=$(GO_SRC)/$(PACKAGE)

GO=go
GO_FMT=$(GO) fmt

.PHONY: build
build: build-overseer build-overlord

.PHONY: build-overseer
build-overseer:
	go build -o ./bin/overseer ./cmd/overseer/main.go

.PHONY: build-overlord
build-overlord:
	go build -o ./bin/overlord ./cmd/overlord/main.go

.PHONY: clean
clean:
	rm -rf ./bin/overseer ./bin/overlord

.PHONY: clean-ci
clean-ci:
	rm -rf ./bin/overseer ./bin/overlord
	rm -rf $(SRC)

.PHONY: bootstrap
bootstrap:
	if [ ! -d "$(SRC)" ]; then mkdir -p "$(GO_SRC)/$(PACKAGE_PATH)" && ln -s "$(PWD)" "$(SRC)" ; fi
	cd "$(SRC)" && go get