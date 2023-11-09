#!make

include .env
export $(grep -v '^#' .env | sed 's/=.*//' | xargs)

current_date := $(shell date +'%Y-%m-%d-%H:%M:%S')

# Run the development server with gin
dev:
	gin -i --appPort 8080 --port 3000 run main.go

# Generate tailwindcss classes for development
tw:
	tailwindcss -i tailwind.base.css -o static/tailwind.css --watch

# Build
build:
	tailwindcss -i tailwind.css -o static/style.css --minify
	docker build -t pengoe .

# Run
run:
	docker run -p 8080:8080 --env-file .env pengoe

# Push migration
push:
	turso db shell $(db) < schema.sql
