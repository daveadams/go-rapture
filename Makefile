gofiles = $(shell find . -name "*.go" )

default: rapture

rapture: $(gofiles)
	go build ./cmd/rapture

test:
	go test ./... -cover

clean:
	rm -f rapture

.PHONY: clean default test
