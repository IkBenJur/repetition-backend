## Migrations
Install in CLI: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

For new migration file run: migrate create -ext sql -dir migrate/migrations -seq FILENAME

Run migration: go run migrate/main.go up/down

## Test
Run test: go test -v ./... 

## Building
sudo docker build --tag TAG_NAME .

## Running container
sudo docker run \
    -e DB_CONN_URL=URL
    -e JWT_SECRET=JWT_STRING