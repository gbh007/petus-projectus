BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

proto:
	protoc -I=. --go_out=. --go-grpc_out=. gate.proto
	protoc -I=. --go_out=. --go-grpc_out=. auth.proto
	protoc -I=. --go_out=. --go-grpc_out=. notification.proto

install-proto:
	sudo apt install protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@lates

build:
	$(BUILD_ENV) go build -o _build/gate ./cmd/gate
	$(BUILD_ENV) go build -o _build/auth ./cmd/authserver
	$(BUILD_ENV) go build -o _build/handler ./cmd/handler
	$(BUILD_ENV) go build -o _build/worker ./cmd/worker
	$(BUILD_ENV) go build -o _build/log ./cmd/log
	$(BUILD_ENV) go build -o _build/notification ./cmd/notification

compose: build
	docker compose up -d --build

desktop:
	go run cmd/desktop-client/main.go