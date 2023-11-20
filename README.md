# Pengoe

Simple solution for splitting money and tracking expenses with friends.

## Stack

- HTMX for interactivity
- turso database
- tailwindcss for styling
- go webserver with html templates

## Commands

- `air` - start dev server in watch mode
  - optional flag: `air -- -log info/warning/error/fatal`
- `make tw` - generate tailwind styles for dev
- `make build` - build docker image
- `make run` - run docker image
- `make push db=<db-name>` - push schema to an empty turso db

### Prerequisites

To run commands, you need to have:

- `make` for running some commands
- `docker` for building and running
- `air` for development server
- `tailwindcss` cli for css class generation
- `bun` for schema migrations with prisma
- `turso` cli for pushing migrations

## Todo

- [x] signup page
  - [x] add username field
  - [x] add svg to background
  - [x] fix tokens
  - [x] update handler for new schema
- [x] signin page
- [x] dashboard page
- [x] loading states
- [x] add services to backend
- [ ] viewtransition api
- [ ] dashboard page
  - [x] handler
  - [x] account selector
  - [x] profile button with signout
  - [ ] show account info
- [x] new account page
  - [x] handler
  - [x] form
- [ ] left panel can be pinned
