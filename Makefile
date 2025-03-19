MIGRATIONS_PATH = ./db/migration
# Ensure .env variables are loaded before executing commands
ifneq (,$(wildcard app.env))
    include app.env
    $(eval export sed 's/=.*//' app.env)
endif

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_SOURCE) --verbose up 

migrate-up1:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_SOURCE) --verbose up 1

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_SOURCE) --verbose down 

migrate-down1:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_SOURCE) --verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/backendproduction-2/db/sqlc Store 

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto  --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090  -r repl

.PHONY: migrate-create migrate-up migrate-down migrate-up1 migrate-down1 db_docs db_schema sqlc test server mock proto evans