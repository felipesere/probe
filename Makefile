build:
	go build .

install: build
	cp probe ~/bin

fmt:
	gofmt -s -w .

lint:
	golangci-lint run --config .golangci.yml
