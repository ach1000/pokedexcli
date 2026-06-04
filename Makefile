.PHONY: build run test

build:
	go build -o pokedexcli .

run:
	go run .

test:
	go test ./...
