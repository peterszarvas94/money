# Pengoe

Simple solution for splitting money and tracking expenses with friends.

## Stack

- HTMX for interactivity
- turso database
- tailwindcss for styling
- go webserver with html templates

## Commands

- `make dev` - start dev server on port 3000
  - optional env: `LOG_LEVEL=INFO/WARNING/ERROR/FATAL`
- `make tw` - generate tailwind styles for dev
- `make build` - build docker image
- `make run` - run docker image
- `make gen` - generate database schema sql
- `make push db=DB_NAME` - push schema to an empty turso db

### Prerequisites

To run make commands, you need to have:

- `make` for running commands
- `docker` for building and running
- `gin` for development server
- `tailwindcss` cli for css class generation
- `bun` for schema migrations with prisma
- `turso` cli for pushing migrations

## Todo

- [x] signup page
  - [ ] add username field
  - [ ] add svg to background
- [ ] signin page
- [ ] dashboard page
- [ ] loading states
- [ ] add models to backend (e.g. User.new, User.login...)
- [ ] viewtransition api
