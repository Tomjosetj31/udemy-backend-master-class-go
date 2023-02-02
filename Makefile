postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://tom:secret@localhost:5432/samplebank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://tom:secret@localhost:5432/samplebank" -verbose down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb postgres dropdb migrateup migratedown sqlc test