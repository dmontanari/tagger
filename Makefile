
VERSION=$(shell git tag --sort=-v:refname | head -n 1)

build:
	cd src/ && go mod download
	cd src/ && go build -ldflags "-X tagger/cmd.versionCode=$(VERSION)"

clean:
	rm -fv src/tagger


