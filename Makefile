# One day we might be able to eat our own dog food but for now this will do.

.PHONY: default restore generate build test clean release publish

default: build

restore:
	go mod download;
	pnpm install;

generate:
	go generate ./resources/;

build: restore generate
	go generate ./resources/;
	go build -o ./dist/gomake ./cmd/gomake/;

test: restore generate
	go test -race -coverprofile ./generator/generator.coverprofile -covermode=atomic ./generator;
	go test -race -coverprofile ./executor/executor.coverprofile -covermode=atomic ./executor;
	go tool cover -html=./generator/generator.coverprofile;
	go tool cover -html=./executor/executor.coverprofile;

release: generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/gomake_linux_amd64 \
		./cmd/gomake/;
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/gomake_darwin_amd64 \
		./cmd/gomake/;
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/gomake_windows_amd64.exe \
		./cmd/gomake/;

clean:
	rm -rf ./dist;
	rm -rf ./node_modules;
	rm -f ./**/*.coverprofile;
	rm -f ./example/.gomake/runner;
	rm -f ./example/.gomake/makefile_generated.go;
