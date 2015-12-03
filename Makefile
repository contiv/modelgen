all: test

clean:
	rm -rf output

godep:
	@if [ -z "`which godep`" ]; then go get -v github.com/kr/godep; fi

vet:
	@(go tool | grep vet) || go get -v golang.org/x/tools/cmd/vet

generator:
	cd generators && bash build.sh >templates.go && gofmt -w -s templates.go

test: godep
	godep go test -v ./...

reflex:
	# go get github.com/cespare/reflex
	reflex -r '.*\.go' -R generators/templates.go make test
