
APP_DIR = cmd/main/main.go

run:
	go run $(APP_DIR)

up:
	docker-compose up --build
