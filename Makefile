# Variables
BACKEND_DIR := ./backend
FRONTEND_DIR := ./frontend

run-backend:
	cd $(BACKEND_DIR) && make run

run-frontend:
	cd $(FRONTEND_DIR) && bun run dev

help:
	@echo "Available targets:"
	@echo "  frontend         Run the frontend application"
	@echo "  backend          Run the backend application"