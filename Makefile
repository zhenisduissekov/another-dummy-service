.PHONY: echo rdc build run test

rdc:
	docker-compose up --remove-orphans --build

build:
	go build -o app ./cmd/dummy-service/main.go

br:
	go build -o app ./cmd/dummy-service/main.go && SERVICE_PORT=:8088 ./app

run:
	go run ./cmd/dummy-service/main.go

test:
	go test ./...

format:
	gofmt -s -w . && goimports -w .


