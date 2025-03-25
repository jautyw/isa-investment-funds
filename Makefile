run:
	go run cmd/server/main.go
mod:
	go mod tidy
lint-install:
 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.8
lint:
	golangci-lint run --timeout=5m
up:
	docker compose up
down:
	docker compose down
docker-prune:
	docker system prune -a -f --volumes
test:
	go test ./...
test-coverage:
	go test ./... -cover
test-coverage-exclusions:
	go test ./... -cover grep -v "/mocks"

generate:
	go generate ./...
