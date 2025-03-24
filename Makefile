run:
	go run cmd/server/main.go
mod:
	go mod tidy
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
