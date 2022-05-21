NAME = silicon-greenhouse-api
BIN = ./bin/$(NAME)
MAIN = main.go

all: build

build: $(MAIN)
	go build -o $(BIN)

arm: $(MAIN)
	env GOOS=linux GOARCH=arm GOARM=5 go build -o $(BIN)-arm .

windows: $(MAIN)
	env GOOS=windows GOARCH=amd64 go build -o $(BIN)-windows.exe .

dev: $(MAIN)
	CompileDaemon -build="make" -command="$(BIN) -port 3000"

clean:
	rm -rf bin/
