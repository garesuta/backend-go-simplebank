MIGRATIONS_PATH = ./db/migration

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) --verbose up 

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) --verbose down 

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: test
test:
	go test -v -cover ./...