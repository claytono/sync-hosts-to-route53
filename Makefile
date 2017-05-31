all: build test lint

FILES=cidrnet.go daemon.go host.go main.go route53.go

build: sync-hosts-to-route53

sync-hosts-to-route53: $(FILES)
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

fmt:
	gofmt -w -s $(FILES)

.PHONY: all clean test build lint fmt
