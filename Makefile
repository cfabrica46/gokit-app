test:
	@echo "Running tests app..."
	go test ./... --cover

cover:
	@echo "Running tests app..."
	go test ./... --coverprofile coverage.out
	go tool cover -func coverage.out

lint:
	@echo "Running golangci-lint app..."
	golangci-lint run

.PHONY: all clean test cover lint
