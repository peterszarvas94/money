#!make

# Include environment variables from .env file
# ENV_FILE := .env
# include $(ENV_FILE)
# export $(shell sed 's/=.*//' $(ENV_FILE))

LOG_LEVEL ?= DEBUG

run:
	killport 8080
	air -c .air.toml -- -log $(LOG_LEVEL)

# Run tests
test:
	air -c .air.test.toml

# Generate templ files
templ:
	air -c .air.templ.toml

# Generate tailwindcss classes for development
tw:
	tailwindcss -i tailwind.base.css -o web/static/tailwind.css --watch

# Docker Build
docker-build:
	tailwindcss -i tailwind.base.css -o web/static/tailwind.css --minify
	docker build -t pengoe .

# Docker Run
docker-run:
	docker run -p 8080:8080 --env-file .env pengoe

# Push migration
push:
	turso db shell $(db) < internal/db/schema.sql

clear:
	turso db shell $(db) < internal/db/clear.sql
