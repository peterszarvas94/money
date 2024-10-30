# Archived

archived, maybe developed furter in the future

# Pengoe

Simple solution for tracking shared accounts.
Split income and log expenses.

## Stack

GOTH

- [go](https://go.dev) server
- [templ](https://github.com/a-h/templ) templates
- [turso](https://turso.tech/) db
- [tailwindcss](https://tailwindcss.com)
- [htmx](https://htmx.org) for reactivity

## Commands

- `make` - start server in watch mode, log errors in logfile
- `make tw` - generate tailwind styles for dev
- `make docker-build` - build docker image
- `make docker-run` - run docker image
- `make push db=<db-name>` - push schema to an empty turso db

### Dependencies

To run commands, you need to have:

- `make` for running some commands
- `docker` for building and running
- `air` for development server
- `tailwindcss` cli for css class generation
- `turso` cli for pushing migrations
- `killport` from my [dotfiles](https://github.com/peterszarvas94/dots/blob/main/.local/bin/killport), because air is buggy

## Todo

### Router

- [ ] switch from custom router to servemux once go 1.22 released

### Auth

- [x] switch to sessions
  - [ ] turso embedded replica
  - [x] session id to session cookie
  - [x] csrf token to hidden `<input />`
  - [ ] maybe csrf token not random but has user info?

### Pages

- [ ] dashboard page
  - [x] handler
  - [x] account selector
  - [x] profile button with signout
  - [ ] show accounts info
  - [ ] maybe some charts
- [ ] account page
  - [x] delete button
  - [x] new event form
  - [x] show events in list
  - [ ] events can have new payment form
  - [x] edit event form
  - [ ] edit payment form

### Components

- [ ] left panel
  - [ ] can be pinned
- [x] event form

### Misc

- [x] switch to uuid
- [ ] add req_id to easch request and log them
- [ ] loading states
  - [x] for protected pages
  - [ ] for button presses
- [ ] toast notifications
- [x] add services to backend
- [ ] viewtransition api
- [ ] better errors
  - [x] central error handling
  - [ ] recover from errors
- [ ] rewrite client-side event bus
