DIR_BIN=bin
EXE=$(DIR_BIN)/github-app-authenticator

.GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
.GIT_HASH=$(shell git rev-list -1 HEAD)

VERSION?=${.GIT_BRANCH}-${.GIT_HASH}
LD_FLAGS=--ldflags "-X main.Version=${VERSION}"

build: $(EXE)

$(EXE): main.go ## Builds the app executable.
	go build $(LD_FLAGS) -o $@

clean: ## Cleans the workspace artifacts.
	rm $(EXE)
