BUILDFLAGS := -mod=vendor

default: all

all: automatix

.PHONY: automatix
automatix:
	$(info * building $@)
	@CGO_ENABLED=0 go build $(BUILDFLAGS) -o "$@"  main.go config.go

.PHONY: check
check:
	$(info * running static code checks)
	@gometalinter ./...
