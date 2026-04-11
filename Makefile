.PHONY: dev build-all docker-build docker-up docker-down docker-logs docker-logs-backend docker-logs-frontend docker-restart docker-clean

dev-backend:
	cd backend && go run cmd/api/main.go

backend-test:
	cd backend && go test ./...

dev-backend-hot:
	cd backend && air

dev-frontend:
	cd frontend && npm run dev

dev:
	make -j 2 dev-backend-hot dev-frontend

install:
	cd backend && go mod tidy
	cd frontend && npm install

# Docker commands
docker-build:
	docker-compose build --no-cache

docker-build-fast:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-restart:
	docker-compose restart

docker-clean:
	docker-compose down -v --rmi local

docker-logs:
	docker-compose logs -f

docker-logs-backend:
	docker-compose logs -f backend

docker-logs-frontend:
	docker-compose logs -f frontend

docker-ps:
	docker-compose ps

# Production deployment
deploy:
	@echo "Deploying to production..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "Deployment complete!"