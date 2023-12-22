# Pengoe

Simple solution for splitting money and tracking expenses with friends.

## Stack

- HTMX for interactivity
- turso database
- tailwindcss for styling
- go webserver with html templates

## Commands

- `make` - start server in watch mode, log errors in logfile
- `make tw` - generate tailwind styles for dev
- `make docker-build` - build docker image
- `make docker-run` - run docker image
- `make push db=<db-name>` - push schema to an empty turso db

### Prerequisites

To run commands, you need to have:

- `make` for running some commands
- `docker` for building and running
- `air` for development server
- `tailwindcss` cli for css class generation
- `bun` for schema migrations with prisma
- `turso` cli for pushing migrations
- `killport` from my [dotfiles](https://github.com/peterszarvas94/dots/blob/main/.local/bin/killport)

## Todo

- [x] signup page
  - [x] add username field
- [x] add svg to background
  - [x] fix tokens
  - [x] update handler for new schema
  - [x] error message if unauthorized
- [x] signin page
  - [x] form
  - [x] error message if unauthorized
- [x] dashboard page
- [ ] loading states
  - [x] for protected pages
  - [ ] for button presses
- [ ] toast notifications
- [x] add services to backend
- [ ] viewtransition api
- [ ] dashboard page
  - [x] handler
  - [x] account selector
  - [x] profile button with signout
  - [ ] show accounts info
  - [ ] maybe some charts
- [x] new account page
  - [x] handler
  - [x] form
  - [x] backend functions
- [ ] left panel can be pinned
- [ ] better errors
  - [x] central error handling
  - [ ] revocer from errors
- [x] topbar 
  - [x] shows accounts correctly
  - [ ] fix empty accountselectitem list 
  - [ ] fix: hidden elements are clickable
- [x] write tests
