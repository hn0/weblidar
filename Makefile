GOPATH=$(shell pwd)
GONAME="lasserver"
PID=/tmp/go-$(GONAME).pid
DATASOURCE=$(GOPATH)"/data/sample.las"

build: clean
	@echo "Building $(GONAME)"
	@GOPATH=$(GOPATH) CGO_CFLAGS_ALLOW="-W*" CGO_CFLAGS="-Wno-deprecated-declarations" go build -tags="cl11" -o bin/$(GONAME) app

run:
	@GOPATH=$(GOPATH) CGO_CFLAGS_ALLOW="-W*" CGO_CFLAGS="-Wno-deprecated-declarations" go run -tags="cl11" src/app/app.go $(DATASOURCE)

test:
	@GOPATH=$(GOPATH) CGO_CFLAGS_ALLOW="-W*" CGO_CFLAGS="-Wno-deprecated-declarations" go test -tags="cl11" clwrapper
	@GOPATH=$(GOPATH) go test -tags="cl11" model

clean:
	@echo "cleaning ..."
	rm -rf bin/$(GONAME)