include .env
export

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

migration:
	goose -dir migrations create $(name) sql

migrate-reset:
	goose -dir migrations postgres "$(DATABASE_URL)" reset