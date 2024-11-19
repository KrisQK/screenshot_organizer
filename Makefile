.PHONY: build

build:
	go build -o screenshot-organizer cmd/main.go
	chmod +x screenshot-organizer

run: build
	./screenshot-organizer 