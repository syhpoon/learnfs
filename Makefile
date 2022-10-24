BUILD=go build -ldflags "-s -w"

.PHONY: fuzz test proto vuln

build:
	@$(BUILD)

proto:
	cd proto && ./gen-pb.sh

test:
	go test -v ./...

fuzz:
	cd pkg/fs && go test -v -fuzz=. -fuzztime 1m

vuln:
	govulncheck ./...

dev:
	rm -f dev.test
	fallocate -l 100m dev.test
	./learnfs mkfs dev.test
