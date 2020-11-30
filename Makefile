build:
	go build .

install: build
	cp probe ~/bin

