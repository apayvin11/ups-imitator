.PHONY: build
build:
	go build -o bin/ups-imitator cmd/ups-imitator/main.go

.PHONY: test
test:
	go test -race -timeout 30s ./...

.DEFAULT_GOAL :=build
