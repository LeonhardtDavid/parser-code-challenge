.PHONY: build
build:
	go build -o bin/parser cmd/parser-challenge/main.go cmd/parser-challenge/config.go

test-coverage:
	@mkdir -p coverage
	go test -coverprofile=coverage/profile.out ./tests/... && \
		go tool cover -html=coverage/profile.out -o coverage/coverage.html
