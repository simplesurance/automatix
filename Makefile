VERSION := $(shell cat ver)
BUILDFLAGS := -mod=vendor
LDFLAGS := -ldflags="-X main.version='$(VERSION)'"
SRC := main.go config.go

default: all

all: automatix

.PHONY: automatix
automatix:
	$(info * building $@ $(VERSION))
	@CGO_ENABLED=0 go build $(BUILDFLAGS) $(LDFLAGS) -o "$@" $(SRC)

.PHONY: check
check:
	$(info * running static code checks)
	@gometalinter ./...

.PHONY: dist/linux_amd64/automatix
dist/linux_amd64/automatix:
	$(info * building $@ $(VERSION))
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o "$@" $(SRC)
	$(info * creating $(@D)/automatix-linux_amd64-$(VERSION).tar.xz)
	@tar $(TARFLAGS) -C $(@D) -cJf $(@D)/automatix-linux_amd64-$(VERSION).tar.xz $(@F)

.PHONY: dirty_worktree_check
dirty_worktree_check:
	@if ! git diff-files --quiet || git ls-files --other --directory --exclude-standard | grep ".*" > /dev/null ; then \
		echo "remove untracked files and changed files in repository before creating a release, see 'git status'"; \
		exit 1; \
		fi

.PHONY: release
release: clean dirty_worktree_check dist/linux_amd64/automatix
	@echo
	@echo next steps:
	@echo - git tag v$(VERSION)
	@echo - git push --tags
	@echo - upload $(shell ls dist/*/*.tar.xz)

.PHONY: clean
clean:
	@rm -rf dist/ automatix
