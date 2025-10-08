## Migrations
Install in CLI: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

For new migration file run: migrate create -ext sql -dir migrate/migrations -seq FILENAME

Run migration: go run migrate/main.go up/down

## Test
Run test: go test -v ./... 