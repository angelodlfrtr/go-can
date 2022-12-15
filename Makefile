GOBIN=go

# Lint code with golint
.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./...

# Clean
.PHONY: clean
clean:
	rm ./coverage.out

.PHONY: test
test:
	go test -v -count 1 -race --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
