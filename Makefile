lint:
	golangci-lint run

vet:
	go vet ./...

ineffassign:
	go get -u github.com/gordonklaus/ineffassign
	ineffassign ./...

test:
	go test ./...

build:
	go build -o bin/main cmd/openvpn-processor/main.go

run:
	go run cmd/openvpn-processor/main.go

cross-compile:
	# 32-Bit Systems
	# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/openvpn-processor/main.go
	# MacOS
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 cmd/openvpn-processor/main.go
	# Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/openvpn-processor/main.go
	# Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/openvpn-processor/main.go
        # 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 cmd/openvpn-processor/main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 cmd/openvpn-processor/main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 cmd/openvpn-processor/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 cmd/openvpn-processor/main.go

upgrade-deps:
	go get -u ./...
	go mod tidy
	go mod vendor

all: test build run
