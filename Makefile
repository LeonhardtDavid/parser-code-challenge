.PHONY: build
build:
	go build -o bin/crawler cmd/crawler/main.go cmd/crawler/config.go

test-coverage:
	@mkdir -p coverage
	go test -coverprofile=coverage/profile.out ./tests/... && \
		go tool cover -html=coverage/profile.out -o coverage/coverage.html
