# Pengoe

Simple solution for tracking shared accounts.
Split income and log expenses.

## Stack

- HTMX for interactivity
- turso database
- tailwindcss for styling
- go webserver with [templ](https://github.com/a-h/templ) templates

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
- `turso` cli for pushing migrations
- `killport` from my [dotfiles](https://github.com/peterszarvas94/dots/blob/main/.local/bin/killport)

## Todo

### Auth

- [ ] switch to sessions
  - [ ] turso embedded replica
  - [ ] session id to session cookie
  - [ ] csrf token to hidden `<input />`

### Pages

- [ ] dashboard page
  - [x] handler
  - [x] account selector
  - [x] profile button with signout
  - [ ] show accounts info
  - [ ] maybe some charts
- [ ] account page
  - [x] delete button
  - [ ] new event form

### Components

- [ ] left panel
  - [ ] can be pinned

### Misc

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
