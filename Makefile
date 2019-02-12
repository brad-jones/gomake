# One day we might be able to eat our own dog food but for now this will do.

.PHONY: default restore generate build test clean release publish build-release-bins

default: build

restore:
	go mod download;
	pnpm install;

generate:
	go generate ./resources/;

build: restore generate
	go generate ./resources/;
	go build -o ./dist/gomake ./cmd/gomake/;

test: clean restore generate
	go test -race -coverprofile ./generator/generator.coverprofile -covermode=atomic ./generator;
	go test -race -coverprofile ./executor/executor.coverprofile -covermode=atomic ./executor;
	go tool cover -html=./generator/generator.coverprofile;
	go tool cover -html=./executor/executor.coverprofile;

build-release-bins:
	rm -rf ./dist;
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/github-downloads/gomake_linux_amd64 \
		./cmd/gomake/;
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/github-downloads/gomake_darwin_amd64 \
		./cmd/gomake/;
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o ./dist/github-downloads/gomake_windows_amd64.exe \
		./cmd/gomake/;

release: generate build-release-bins
	mkdir -p ./dist/homebrew-tap;
	VERSION=$(VERSION) HASH=$(shell sha256sum ./dist/github-downloads/gomake_darwin_amd64 | head -c 64) \
		envsubst < ./brew.rb > ./dist/homebrew-tap/gomake.rb;
	mkdir -p ./dist/scoop-bucket;
	VERSION=$(VERSION) HASH=$(shell sha256sum ./dist/github-downloads/gomake_windows_amd64.exe | head -c 64) \
		envsubst < ./scoop.json > ./dist/scoop-bucket/gomake.json;

publish:
	git clone --progress https://${GITHUB_TOKEN}@github.com/brad-jones/homebrew-tap.git /tmp/homebrew-tap;
	rm -f /tmp/homebrew-tap/Formula/gomake.rb;
	cp ./dist/homebrew-tap/gomake.rb /tmp/homebrew-tap/Formula/gomake.rb;
	cd /tmp/homebrew-tap && \
	git add -A && \
	git commit -m "chore(gomake): release new version $(VERSION)" && \
	git push origin master;
	git clone --progress https://${GITHUB_TOKEN}@github.com/brad-jones/scoop-bucket.git /tmp/scoop-bucket;
	rm -f /tmp/scoop-bucket/gomake.json;
	cp ./dist/scoop-bucket/gomake.json /tmp/scoop-bucket/gomake.json;
	cd /tmp/scoop-bucket && \
	git add -A && \
	git commit -m "chore(gomake): release new version $(VERSION)" && \
	git push origin master;

clean:
	rm -rf ./dist;
	rm -rf ./node_modules;
	rm -f ./**/*.coverprofile;
	rm -f ./example/.gomake/runner;
	rm -f ./example/.gomake/makefile_generated.go;
