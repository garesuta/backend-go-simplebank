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

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/backendproduction-2/db/sqlc Store 


.PHONY: migrate-create migrate-up migrate-down migrate-up1 migrate-down1 sqlc test server mock