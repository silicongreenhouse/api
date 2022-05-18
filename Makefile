NAME = silicon-greenhouse-api
BIN = ./bin/$(NAME)
MAIN = main.go

all: build

build: $(MAIN)
	go build -o $(BIN)

dev: $(MAIN)
	CompileDaemon -build="make" -command="$(BIN)"

clean:
	rm -rf bin/