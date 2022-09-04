BUILD=go build -ldflags "-s -w"
PKG=learnfs/cmd

.PHONY: fuzz test proto

build:
	@$(BUILD) -o learnfs $(PKG)

proto:
	cd proto && ./gen-pb.sh

test:
	go test -v learnfs/...

dev:
	rm -f dev.test
	fallocate -l 10m dev.test
	./learnfs mkfs dev.test
