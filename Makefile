lint:
	@golangci-lint run ./...

test:
	@go test -v -cover ./...

coverage:
	@go test -coverprofile=coverage.out ./...

coverage-html:
	@go tool cover -html=coverage.out

ci: lint test coverage
