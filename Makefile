VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
GOBIN := $(GOPATH)/bin

export GO111MODULE = on

build_tags = netgo

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
empty = $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(empty),$(comma),$(build_tags))

ldflags = -X github.com/bandprotocol/falcon/main.Name=falcon \
	-X github.com/bandprotocol/falcon/cmd.AppName=falcon \
	-X github.com/bandprotocol/falcon/cmd.Commit=$(COMMIT) \
	-X github.com/bandprotocol/falcon/cmd.Version=$(VERSION) \
	-X "github.com/bandprotocol/falcon/cmd.BuildTags=$(build_tags_comma_sep)"

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)'

PACKAGE_NAME         := github.com/bandprotocol/falcon
GOLANG_CROSS_VERSION ?= latest

all: install

install: go.sum
	@echo "installing falcon binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/falcon main.go

.PHONY: release
#? release: Run goreleaser to build and release cross-platform falcon binary version
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean

PROTOC        := protoc
PROTO_DIR     := proto
OUT_DIR       := .


.PHONY: proto
#? proto: Generate Go code from protobuf definitions
proto:
	@echo "Generating protobuf code..."
	@$(PROTOC) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/**/**/*.proto

mocks:
	@go install go.uber.org/mock/mockgen@latest
	sh ./scripts/mockgen.sh

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify
	touch go.sum

test:
	@go test -mod=readonly ./...

test-coverage:
	@go test -mod=readonly -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
