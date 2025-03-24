run:
	go run cmd/server/main.go
mod:
	go mod tidy
lint-install:
 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
lint:
	golangci-lint run --timeout=5m --skip-dirs "/pkg/mod/"
run-docker:
	docker compose up
stop-docker:
	docker compose down
docker-prune:
	docker system prune -a -f --volumes
test:
	go test ./...
test-coverage:
	go test ./... -cover
generate:
	go generate ./...
