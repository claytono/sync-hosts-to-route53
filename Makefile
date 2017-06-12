all: build test lint

VERSION=$(shell git describe --dirty)
FILES=cidrnet.go daemon.go host.go main.go route53.go
BINS=sync-hosts-to-route53-linux-mips64 \
	sync-hosts-to-route53-linux-mips \
	sync-hosts-to-route53-linux-arm \
	sync-hosts-to-route53-linux-arm64 \
	sync-hosts-to-route53-linux-386 \
	sync-hosts-to-route53-linux-amd64 \

BUILDFLAGS=-v -ldflags "-X main.version=$(VERSION)"

build: sync-hosts-to-route53

build-all-arch: $(BINS)

sync-hosts-to-route53: $(FILES)
	go build $(BUILDFLAGS)

sync-hosts-to-route53-linux-mips64: $(FILES)
	GOOS=linux GOARCH=mips64 go build $(BUILDFLAGS) -o $@

sync-hosts-to-route53-linux-mips: $(FILES)
	GOOS=linux GOARCH=mips go build $(BUILDFLAGS) -o $@

sync-hosts-to-route53-linux-arm: $(FILES)
	GOOS=linux GOARCH=arm go build $(BUILDFLAGS) -o $@

sync-hosts-to-route53-linux-arm64: $(FILES)
	GOOS=linux GOARCH=arm64 go build $(BUILDFLAGS) -o $@

sync-hosts-to-route53-linux-386: $(FILES)
	GOOS=linux GOARCH=386 go build $(BUILDFLAGS) -o $@

sync-hosts-to-route53-linux-amd64: $(FILES)
	GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o $@

test:
	go test -v $(glide novendor)

lint:
	golint

coverage-html:
	go test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	go clean
	rm -f ${BINS}

fmt:
	gofmt -w -s $(FILES)

.PHONY: all clean test build lint fmt
