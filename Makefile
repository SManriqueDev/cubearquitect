.PHONY: dev build-all

dev-backend:
	cd backend && go run cmd/api/main.go

dev-frontend:
	cd frontend && npm run dev

dev:
	make -j 2 dev-backend dev-frontend

install:
	cd backend && go mod tidy
	cd frontend && npm install