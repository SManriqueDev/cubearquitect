.PHONY: dev build-all docker-build docker-up docker-down

dev-backend:
	cd backend && go run cmd/api/main.go

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
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-logs-backend:
	docker-compose logs -f backend

docker-logs-frontend:
	docker-compose logs -f frontend