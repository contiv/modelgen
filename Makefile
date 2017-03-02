
.PHONY: all build generator godep test

all: build test

# build runs the generator script and compiles and installs the code
build: generator
	go install -v

# generator runs a script which dynamically generates a .go file from the
# .tmpl files in the "templates" directory
generator:
	cd generators && bash build.sh >templates.go && gofmt -w -s templates.go

# godep updates Godeps/Godeps.json
godep:
	godep save ./...

# test runs all the tests
test:
	go test -v ./...
