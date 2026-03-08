.PHONY: up down test test-integration

up:
	docker compose up --build -d

down:
	docker compose down -v

test:
	go test ./...

test-integration:
	go test -tags=integration ./test/integration/...
