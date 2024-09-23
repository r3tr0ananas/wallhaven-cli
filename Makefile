build:
	go build -o wallhaven ./src/

install:
	cp wallhaven ~/.local/bin

run:
	go run ./src/ $(ARGS)