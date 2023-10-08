.PHONY: setup
setup:
	go install go.uber.org/mock/mockgen@latest

.PHONY: gen
gen:
	go generate ./...

TARGET := .

.PHONY: test
test:
	go test -race -shuffle=on -run $(TARGET) ./...

.PHONY: server/run
server/run:
	go run cmd/server/main.go

.PHONY: client/run
client/run:
	go run cmd/client/main.go
