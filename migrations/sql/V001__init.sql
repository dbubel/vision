CREATE EXTENSION if not exists vector;
CREATE TABLE IF NOT EXISTS workspaces
(
    id          SERIAL PRIMARY KEY,
    name        TEXT      NOT NULL UNIQUE,
    description TEXT      NOT NULL DEFAULT '',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workspace_vector_mappings
(
    id           SERIAL PRIMARY KEY,
    name         text UNIQUE     not null,
    workspace_id INTEGER REFERENCES workspaces (id) ON DELETE CASCADE,
    dimension    INTEGER   NOT NULL,
    vector_table TEXT      NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp
);
