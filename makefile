CMD=/usr/local/go/bin/go
BIN_PATH=bin
SRC_PATH=src

all: clean build install

build:
	$(CMD) build -o $(BIN_PATH)/zind $(SRC_PATH)/*

install:
	cp bin/zind /usr/bin/zind
	mkdir -p /var/lib/zind/images
	mkdir -p /var/lib/zind/volumes
	mkdir -p /var/lib/zind/containers

uninstall:
	rm -rf bin/zind
	rm -rf /usr/bin/zind

clean: uninstall
