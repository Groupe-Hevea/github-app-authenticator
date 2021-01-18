VERSION?=$(shell git rev-list -1 HEAD)
DIR_BIN=bin
EXE=github-app-authenticator

LD_FLAGS=-ldflags "-X main.Version=$VERSION"

build: $(DIR_BIN)/$(BIN)

$(DIR_BIN)/$(BIN): ## Builds the app executable
	go build $(LD_FLAGS) -o $@
