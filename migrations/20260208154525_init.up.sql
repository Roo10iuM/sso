CREATE TABLE
    IF NOT EXISTS users (
        id BLOB PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        email TEXT UNIQUE,
        pass_hash BLOB NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_username ON users (username);

CREATE TABLE
    IF NOT EXISTS apps (
        id BLOB PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        secret TEXT NOT NULL UNIQUE
    );