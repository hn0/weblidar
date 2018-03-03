GOPATH=$(shell pwd)
GONAME="lasserver"
PID=/tmp/go-$(GONAME).pid
DATASOURCE=$(GOPATH)"/data/sample.las"
CFLAGS="-I/usr/include/CL"
LDFLAGS="-L/usr/lib/clc"

build:
	@echo "Building $(GONAME)"
	@GOPATH=$(GOPATH) go build -v -o bin/$(GONAME) app

run:
	@GOPATH=$(GOPATH) LDFLAGS=${LDFLAGS} CGO_CFLAGS=${CFLAGS} go run src/app/app.go $(DATASOURCE)

test:
	@GOPATH=$(GOPATH) CGO_CFLAGS=${CFLAGS} go test clwrapper
	@GOPATH=$(GOPATH) go test model