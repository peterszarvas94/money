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

### Pages

- [x] signup page
  - [x] add username field
- [x] signin page
  - [x] form
  - [x] error message if unauthorized
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
- [ ] account page
  - [x] delete button
  - [ ] new event form

### Components

- [x] topbar 
  - [x] shows accounts correctly
  - [x] fix empty accountselectitem list 
  - [x] fix: hidden elements are clickable
- [ ] left panel
  - [ ] can be pinned

### Misc

- [x] add svg to background
  - [x] fix tokens
  - [x] update handler for new schema
  - [x] error message if unauthorized
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
- [x] tests
