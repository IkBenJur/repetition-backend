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
    -e DATABASE_URL=URL
    -e JWT_SECRET=JWT_STRING

### Running container local
When using Linux add host to run command: --add-host=host.docker.internal:host-gateway 
Use: postgres://username:password@host.docker.internal:5432/database