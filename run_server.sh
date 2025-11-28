#!/bin/sh
set -e

echo "Running migrations"
go run migrate/main.go up

echo "Starting app"
./repetition-backend