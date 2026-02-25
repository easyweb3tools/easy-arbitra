APP=easy-arbitra-server

.PHONY: tidy
 tidy:
	go mod tidy

.PHONY: run
 run:
	go run ./cmd/server

.PHONY: seed
 seed:
	go run ./cmd/seed

.PHONY: build
 build:
	go build -o bin/$(APP) ./cmd/server

.PHONY: test
 test:
	go test ./...
