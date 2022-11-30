PROJECT ?= $(shell basename $(CURDIR))
MODULE  ?= $(shell go list -m)

GO      ?= GO111MODULE=on go
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
BIDTIME ?= $(shell date +%FT%T%z)

BITTAGS := viper_logger
LDFLAGS := -s -w
LDFLAGS += -X "github.com/starudream/go-lib/constant.VERSION=$(VERSION)"
LDFLAGS += -X "github.com/starudream/go-lib/constant.BIDTIME=$(BIDTIME)"
LDFLAGS += -X "github.com/starudream/go-lib/constant.PREFIX=ST"

.PHONY: bin

bin:
	@$(MAKE) bin-client
	@$(MAKE) bin-server

bin-client:
	@$(MAKE) tidy
	CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/stc $(MODULE)/cmd/client

bin-server:
	@$(MAKE) tidy
	CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/sts $(MODULE)/cmd/server

run-client:
	@$(MAKE) tidy
	CGO_ENABLED=1 $(GO) run -race -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' $(MODULE)/cmd/client

run-server:
	@$(MAKE) tidy
	CGO_ENABLED=1 $(GO) run -race -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' $(MODULE)/cmd/server

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/*

upx:
	upx bin/*

lint:
	golangci-lint run --skip-dirs-use-default
