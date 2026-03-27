.PHONY: test lint fmt bench

test:
	go test ./... -v -race

bench:
	go test ./... -bench=. -benchmem

lint:
	golangci-lint run ./...

fmt:
	gofumpt -w .
	goimports -w .
