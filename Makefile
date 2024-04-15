.PHONY: run-tests
run-tests: 
		@go test -v -cover `go list ./...` 