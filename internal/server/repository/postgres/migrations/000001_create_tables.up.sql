CREATE TABLE IF NOT EXISTS keeper.users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at    TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS keeper.secrets
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    user_id UUID         NOT NULL REFERENCES keeper.users (id) ON DELETE CASCADE,
    hash    VARCHAR(255) NOT NULL,
    data    BYTEA        NOT NULL
);

ALTER TABLE keeper.secrets
    ADD CONSTRAINT secrets_name_user_id_key UNIQUE (name, user_id);

CREATE INDEX IF NOT EXISTS idx_secrets_user_id ON keeper.secrets (user_id);

CREATE TYPE keeper.secret_status AS ENUM ('ACTIVE', 'DELETED');

CREATE TABLE IF NOT EXISTS keeper.secrets_statuses
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255)         NOT NULL,
    user_id       UUID                 NOT NULL REFERENCES keeper.users (id) ON DELETE CASCADE,
    last_modified TIMESTAMP            NOT NULL,
    status        keeper.secret_status NOT NULL
);

ALTER TABLE keeper.secrets_statuses
    ADD CONSTRAINT secrets_statuses_name_user_id_key UNIQUE (name, user_id);

CREATE INDEX IF NOT EXISTS idx_secrets_user_id ON keeper.secrets_statuses (user_id);
