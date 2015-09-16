fmt: format

format:
	go fmt ./...

bin: fmt
	go build -o bin/envar -ldflags "-X main.Version $$(git describe --tags)" *.go
