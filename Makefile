all: install

install:
	go install ./...

.PHONY: examples
examples:
	cd examples/ && ./compile.sh

