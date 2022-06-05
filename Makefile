.PHONY: test dep

test:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

dep:
	go mod tidy
	go mod download
	go mod verify
