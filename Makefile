GOBIN=go

# Lint code with golint
.PHONY: lint
lint:
	find . -name "*.go" | xargs misspell -error
	golint -set_exit_status ./...
	go vet ./...
	staticcheck ./...

# Clean
.PHONY: clean
clean:
	rm -Rf ./coverage.out

.PHONY: test
test:
	go test -v -count 1 -race --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
