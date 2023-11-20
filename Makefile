#!make

# Generate tailwindcss classes for development
tw:
	tailwindcss -i tailwind.base.css -o web/static/tailwind.css --watch

# Docker Build
build:
	tailwindcss -i tailwind.base.css -o web/static/tailwind.css --minify
	docker build -t pengoe .

# Docker Run
run:
	docker run -p 8080:8080 --env-file .env pengoe

# Push migration
push:
	turso db shell $(db) < schema.sql
