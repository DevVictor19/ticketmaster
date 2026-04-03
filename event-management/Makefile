.PHONY: migrate-create migrate-up migrate-down migrate-status

migrate-create:
	goose create $(name) sql

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

.PHONY: run-dev 

run-dev:
	air