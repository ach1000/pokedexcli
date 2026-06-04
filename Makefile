.PHONY: build run test clean

build:
	go build -o pokedexcli .

run:
	go run .

test:
	go test ./...

clean:
	rm -f pokedexcli repl.log
