-- Turso SQLite3 Database Schema
CREATE TABLE
  user (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
  );

CREATE TABLE
  account (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    currency TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
  );

CREATE TABLE
  access (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    role TEXT CHECK (role IN ('admin', 'viewer')) NOT NULL,
    user_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE ON UPDATE CASCADE
  );

CREATE TABLE
  recipient (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    access_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (access_id) REFERENCES access (id) ON DELETE CASCADE ON UPDATE CASCADE
  );

CREATE TABLE
  event (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    income INTEGER NOT NULL,
    reserved INTEGER NOT NULL,
    delivered_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE ON UPDATE CASCADE
  );

CREATE TABLE
  payment (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    factor INTEGER NOT NULL,
    extra INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    paid INTEGER NOT NULL CHECK (paid IN (0, 1)),
    paid_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES event (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES recipient (id) ON DELETE CASCADE ON UPDATE CASCADE
  );

CREATE TABLE
  session (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    valid_until DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
  );
