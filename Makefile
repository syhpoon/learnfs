BUILD=go build -ldflags "-s -w"
PKG=github.com/syhpoon/learnfs/cmd

.PHONY: fuzz test proto vuln

build:
	@$(BUILD) -o learnfs $(PKG)

proto:
	cd proto && ./gen-pb.sh

test:
	go test -v ./...

vuln:
	govulncheck ./...

dev:
	rm -f dev.test
	fallocate -l 10m dev.test
	./learnfs mkfs dev.test
