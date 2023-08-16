postgres:
	@echo "Setting up postgres..."
	docker run --name broempSignalDB -e POSTGRES_USER=broempSignal -e POSTGRES_PASSWORD=broempSignal -p 5432:5432 -d postgres:alpine

createdb:
	@echo "Creating database..."
	docker exec -it broempSignalDB createdb --username=broempSignal --owner=broempSignal broempSignal

dropdb:
	@echo "Dropping database..."
	docker rm -f broempSignalDB

migrateup:
	migrate -path db/migration -database "postgresql://broempSignal:broempSignal@$(DB_URL):5432/broempSignal?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://broempSignal:broempSignal@$(DB_URL):5432/broempSignal?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server
