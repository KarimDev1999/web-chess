.PHONY: fmt vet test build clean migrate-up migrate-down migrate-create

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

test:
	go test -race -cover ./...

build:
	go build ./...

clean:
	go clean -cache -testcache -modcache

lint: fmt vet test

migrate-up:
	migrate -path ./migrations -database "postgres://chess:chess@localhost:5432/chess?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://chess:chess@localhost:5432/chess?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name
