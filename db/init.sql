CREATE TABLE IF NOT EXISTS repository_tb
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    author     VARCHAR(255) NOT NULL,
    parse_status VARCHAR(20)  NOT NULL DEFAULT 'PENDING',
    api_url    TEXT         NOT NULL,
    gh_url     TEXT         NOT NULL,
    clone_url  TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, author)
);

CREATE TABLE IF NOT EXISTS parser_state_tb
(
    id                   SERIAL PRIMARY KEY,
    last_parsed_repo_id  INTEGER REFERENCES repository_tb(id),
    last_parsed_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS parser_links_tb
(
    id           SERIAL PRIMARY KEY,
    repo_id      INTEGER NOT NULL REFERENCES repository_tb(id),
    url          TEXT NOT NULL,
    expanded_url TEXT,
    file         TEXT NOT NULL,
    line_number  INTEGER NOT NULL,
    path         TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (repo_id, url, file, line_number)
);
