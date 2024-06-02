.PHONY: run-tests
run-tests: 
		@go test -v -cover `go list ./...`
		
.PHONY: lint
lint:
	@`go env GOPATH`/bin/golangci-lint run

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source ./$(util)/$(subutil).go -destination ./tests/mock/$(util)/$(subutil).go