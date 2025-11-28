# syntax=docker/dockerfile:1

# Get the base image
FROM golang:1.25

# Create base workdit
WORKDIR /app

# Copy over mod and sum files
COPY go.mod go.sum ./

# Install depencies
RUN go mod download

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy over source code
COPY . ./

# Create binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /repetition-backend

# Copy bash running script
COPY run_server.sh ./

# Set run permission
RUN chmod +x run_server.sh

# Run executable
CMD ["./run_server.sh"]