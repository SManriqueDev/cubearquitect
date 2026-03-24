.PHONY: dev build-all

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