build:
	go build -o wallhaven ./wallhaven/

install:
	cp wallhaven ~/.local/bin

run:
	go run ./wallhaven/ $(ARGS)