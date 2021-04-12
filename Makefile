BINARY            = tswitch
VERSION          ?= $(shell git describe --tags --dirty --always | sed -e 's/^v//')
GOCMD             = go
MOCKGEN_VERSION   = 1.5.0
GOFLAGS          ?= $(GOFLAGS:)
LDFLAGS          := "-X main.version=$(VERSION)"
IS_SNAPSHOT       = $(if $(findstring -, $(VERSION)),true,false)
MAJOR_VERSION     = $(word 1, $(subst ., ,$(VERSION)))
MINOR_VERSION     = $(word 2, $(subst ., ,$(VERSION)))
PATCH_VERSION     = $(word 3, $(subst ., ,$(word 1,$(subst -, , $(VERSION)))))
NEW_VERSION      ?= $(MAJOR_VERSION).$(MINOR_VERSION).$(shell echo $$(( $(PATCH_VERSION) + 1)) )

all: build

${BINARY}:
	"$(GOCMD)" build ${GOFLAGS} -ldflags ${LDFLAGS} -o dist/"${BINARY}"

build: ${BINARY}

.PHONY: mocks
mocks:
	@go install github.com/golang/mock/mockgen@v${MOCKGEN_VERSION}
	mockgen -destination=mocks/db.go -package=mocks -source=pkg/db/db.go
	mockgen -destination=mocks/installer.go -package=mocks -source=pkg/lib/installer.go

.PHONY: test
test: build
	"$(GOCMD)" test -race -v ./...

.PHONY: clean
clean:
	"$(GOCMD)" clean -i
	@rm -fr dist/*

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: promote
promote:
	@git fetch --tags
	@echo "VERSION:$(VERSION) IS_SNAPSHOT:$(IS_SNAPSHOT) NEW_VERSION:$(NEW_VERSION)"
ifeq (false,$(IS_SNAPSHOT))
	@echo "Unable to promote a non-snapshot"
	@exit 1
endif
ifneq ($(shell git status -s),)
	@echo "Unable to promote a dirty workspace"
	@exit 1
endif
	git tag -a -m "releasing v$(NEW_VERSION)" v$(NEW_VERSION)
	git push origin v$(NEW_VERSION)
