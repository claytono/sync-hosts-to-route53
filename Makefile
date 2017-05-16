all: build test lint

build: sync-hosts-to-route53

sync-hosts-to-route53:
	go build -v

test:
	go test -v $(glide novendor)

lint:
	golint

clean:
	go clean

.PHONY: all clean test build lint
