all: build test lint

build: sync-hosts-to-route53

sync-hosts-to-route53: cidrnet.go host.go main.go route53.go
	go build -v

test:
	go test -v $(glide novendor)

lint:
	golint

coverage-html:
	go test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	go clean

.PHONY: all clean test build lint
