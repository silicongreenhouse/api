NAME = silicon-greenhouse-api
BIN = ./bin/$(NAME)
MAIN = main.go

all: build

build: $(MAIN)
	make linux
	make arm
	make windows
	make macos

arm: $(MAIN)
	env GOOS=linux GOARCH=arm GOARM=5 go build -o $(BIN)-arm .

windows: $(MAIN)
	env GOOS=windows GOARCH=amd64 go build -o $(BIN)-windows.exe .

macos: $(MAIN)
	env GOOS=darwin GOARCH=amd64 go build -o $(BIN)-macos .

linux: $(MAIN)
	env GOOS=linux GOARCH=amd64 go build -o $(BIN)-linux .

dev: $(MAIN)
	CompileDaemon -build="go build -o $(BIN)-env" -command="$(BIN)-env -port 3000"

clean:
	rm -rf bin/
