install:
	go mod download

create_token:
	go run main.go create token

read_token:
	go run main.go read database

all: install create_token read_token

compile:
	echo "Compiling for every OS and Platform"
	GOOS=windows go build -o bin/main-amd64.exe main.go
	GOOS=darwin go build -o bin/main-amd64-darwin main.go
	GOOS=linux  go build -o bin/main-linux-arm64 main.go
