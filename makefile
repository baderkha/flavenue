start-db:
	cd docker && docker-compose up -d
start:
	go run cmd/local/main.go
fake-data:
	go run scripts/fake_data/main.go
migrate:
	go run scripts/db/migrate/migrate.go