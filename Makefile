.PHONY: run-tests
run-tests: 
		@go test -v -cover `go list ./...`
		
.PHONY: lint
lint:
	@`go env GOPATH`/bin/golangci-lint run