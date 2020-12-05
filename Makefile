postgres:
	docker create --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:12-alpine

start_postgres:
	docker start postgres12

stop_postgres:
	docker stop postgres12

remove_postgres:
	docker container rm postgres12

create_migration:
	migrate create -ext sql -dir db/migration -seq init_schema

create_db:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

drop_db:
	docker exec -it postgres12 dropdb simple_bank

migrate_up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...