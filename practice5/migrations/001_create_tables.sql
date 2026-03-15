-- 001_create_tables.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) UNIQUE NOT NULL,
    gender     VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    friend_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT no_self_friendship CHECK (user_id <> friend_id)
);
