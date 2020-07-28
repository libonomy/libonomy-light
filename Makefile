BINARY := go-libonomy
COMMIT = $(shell git rev-parse HEAD)
SHA = $(shell git rev-parse --short HEAD)
CURR_DIR = $(shell pwd)
CURR_DIR_WIN = $(shell cd)
BIN_DIR = $(CURR_DIR)/build
BIN_DIR_WIN = $(CURR_DIR_WIN)/build
export GO111MODULE = on

PKGS = $(shell go list ./...)

PLATFORMS := windows linux darwin
os = $(word 1, $@)



install:
ifeq ($(OS),Windows_NT) 
	setup_env.bat
else
	./setup_env.sh
endif
.PHONY: install


p2p-build:
ifeq ($(OS),WINDOWS_NT)
	cd cmd/p2p ; go build -o $(BIN_DIR_WIN)/p2p-simulate.exe; cd ..
else
	cd cmd/p2p ; go build -o $(BIN_DIR)/p2p-simulate; cd ..
endif
.PHONY: p2p


