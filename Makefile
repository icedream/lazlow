VERSION = $(shell git describe --always --tags --dirty)

SOURCES = $(wildcard *.go)

LAZLOW_SOURCES = $(wildcard cmd/lazlow/*.go)

LAZLOW_OUTPUT = lazlow

.PHONY: default
default: all

.PHONY: all
all: $(LAZLOW_OUTPUT)

$(LAZLOW_OUTPUT): $(SOURCES) $(LAZLOW_SOURCES)
	go build -v -ldflags="-X 'main.Version=$(VERSION)' -s -w" -o $(LAZLOW_OUTPUT) ./cmd/lazlow
