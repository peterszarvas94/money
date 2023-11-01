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
	docker build -t go-hmtx .

# Run
run:
	docker run -p 8080:8080 --env-file .env go-hmtx

# Set latest migration
MIGRATIONS_DIR := db/migrations
MIGRATIONS_EXISTS := $(wildcard $(MIGRATIONS_DIR)/*.sql)
ifneq ($(MIGRATIONS_EXISTS),)
	latest := $(shell ls -t $(MIGRATIONS_DIR)/*.sql | head -1)
endif

# Generate migration
gen:
	cd db && \
	mkdir -p migrations && \
	bunx prisma migrate diff --from-empty --to-schema-datamodel schema.prisma --script > migrations/prisma-$(current_date).sql

# Push migration
push:
	turso db shell $(db) < $(latest)
