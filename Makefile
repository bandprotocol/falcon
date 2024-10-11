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

ldflags = -X github.com/bandprotocol/falcon/version.Name=falcon \
	-X github.com/bandprotocol/falcon/version.AppName=falcon \
	-X github.com/bandprotocol/falcon/version.Commit=$(COMMIT) \
	-X github.com/bandprotocol/falcon/version.Version=$(VERSION) \
	-X "github.com/bandprotocol/falcon/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)'

all: install

install: go.sum
	@echo "installing falcon binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(GOBIN)/falcon main.go

release: go.sum
	env GOOS=linux GOARCH=amd64 \
		go build -mod=readonly -o ./build/falcon_linux_amd64 $(BUILD_FLAGS) main.go
	env GOOS=darwin GOARCH=amd64 \
		go build -mod=readonly -o ./build/falcon_darwin_amd64 $(BUILD_FLAGS) main.go
	env GOOS=windows GOARCH=amd64 \
		go build -mod=readonly -o ./build/falcon_windows_amd64 $(BUILD_FLAGS) main.go

mocks:
	@go install go.uber.org/mock/mockgen@latest
	sh ./scripts/mockgen.sh

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify
	touch go.sum

test:
	@go test -mod=readonly ./...
