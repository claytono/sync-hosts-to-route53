all: build test

build: sync-hosts-to-route53

sync-hosts-to-route53:
	go build -v

test:
	go test -v $(glide novendor)

clean:
	go clean

.PHONY: all clean test build