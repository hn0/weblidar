GOPATH=$(shell pwd)
GONAME="lasserver"
PID=/tmp/go-$(GONAME).pid
DATASOURCE=$(GOPATH)"/data/sample.las"

build: clean
	@echo "Building $(GONAME)"
	@GOPATH=$(GOPATH) go build -tags="cl11" -gcflags="-Wno-deprecated-declarations" -o bin/$(GONAME) app

run:
	@GOPATH=$(GOPATH) go run -tags="cl11" src/app/app.go $(DATASOURCE)

test:
	@GOPATH=$(GOPATH) go test -tags="cl11" clwrapper
	@GOPATH=$(GOPATH) go test -tags="cl11" model

clean:
	@echo "cleaning ..."
	rm -rf bin/$(GONAME)