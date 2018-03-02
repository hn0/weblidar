GOPATH=$(shell pwd)
GONAME="lasserver"
PID=/tmp/go-$(GONAME).pid
DATASOURCE=$(GOPATH)"/data/tmp"

build:
	@echo "Building $(GONAME)"
	@GOPATH=$(GOPATH) go build -v -o bin/$(GONAME) app

run:
	@GOPATH=$(GOPATH) go run src/app/app.go $(DATASOURCE)

test:
	@GOPATH=$(GOPATH) go test template