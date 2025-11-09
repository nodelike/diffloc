.DEFAULT_GOAL := help
.PHONY: help build install run test clean release-test

help:
	@echo "diffloc - Available commands:"
	@echo "  make build        - Build binary to bin/diffloc"
	@echo "  make install      - Install to GOPATH"
	@echo "  make run          - Run without building"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make release-test - Test release build locally"

build:
	go build -o bin/diffloc ./cmd/diffloc

install:
	go install ./cmd/diffloc

run:
	go run ./cmd/diffloc

clean:
	rm -rf bin/ dist/

release-test:
	goreleaser release --snapshot --clean

