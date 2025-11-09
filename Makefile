build:
	go build -o bin/diffloc ./cmd/diffloc

install:
	go install ./cmd/diffloc

run:
	go run ./cmd/diffloc

test:
	go test ./...

clean:
	rm -rf bin/ dist/

release-test:
	goreleaser release --snapshot --clean

