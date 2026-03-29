
VERSION=$(shell git tag --sort=-v:refname | head -n 1)

build:
	go mod download
	go build -ldflags "-X tagger/cmd.versionCode=$(VERSION)"

clean:
	rm -fv src/tagger


