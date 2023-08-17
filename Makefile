BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

proto:
	protoc -I=. --go_out=. --go-grpc_out=. services/gate/gate.proto
	protoc -I=. --go_out=. --go-grpc_out=. services/auth/auth.proto
	protoc -I=. --go_out=. --go-grpc_out=. services/notification/notification.proto
	protoc -I=. --go_out=. --go-grpc_out=. services/log/log.proto

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

up: build
	docker compose up -d --build

logs:
	docker compose logs -f auth gate handler log notification worker

down:
	docker compose down --remove-orphans

desktop:
	go run cmd/desktop-client/main.go