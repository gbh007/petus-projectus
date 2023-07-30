proto:
	protoc -I=. --go_out=. --go-grpc_out=. gate.proto
	protoc -I=. --go_out=. --go-grpc_out=. auth.proto

install-proto:
	sudo apt install protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@lates

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o _build/gate ./cmd/gate
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o _build/auth ./cmd/authserver

compose: build
	docker-compose up -d --build