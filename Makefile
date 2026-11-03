.PHONY: dev build test up down load-test clean

# Run backend locally
dev-api:
	cd apps/api && go run cmd/server/main.go

# Run frontend locally
dev-web:
	cd apps/web && npm run dev

# Run unit tests
test:
	cd apps/api && go test -v -race ./...

# Start full multi-container docker stack
up:
	docker compose -f deployments/docker-compose.yml up -d --build

# Stop docker stack
down:
	docker compose -f deployments/docker-compose.yml down -v

# Run WebSocket k6 load test
load-test:
	k6 run scripts/ws_load_test.js

clean:
	rm -rf apps/api/bin apps/web/.next
