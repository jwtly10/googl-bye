CREATE DATABASE "googl-bye-db"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c googl-bye-db;

CREATE TABLE IF NOT EXISTS public.repository_tb (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    language TEXT DEFAULT '',
    stars INTEGER DEFAULT 0,
    forks INTEGER DEFAULT 0,
    size INTEGER DEFAULT 0,
    last_push TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    state VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    api_url TEXT NOT NULL,
    gh_url TEXT NOT NULL,
    clone_url TEXT NOT NULL,
    error_msg TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, author)
);

CREATE TABLE IF NOT EXISTS public.search_params_history_tb (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    query TEXT DEFAULT '',
    opts JSON,
    start_page INTEGER NOT NULL DEFAULT 0,
    current_page INTEGER NOT NULL DEFAULT 0,
    pages_to_process INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS public.parser_state_tb (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    last_parsed_repo_id INTEGER REFERENCES public.repository_tb(id),
    last_parsed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS public.parser_links_tb (
    id SERIAL PRIMARY KEY,
    repo_id INTEGER NOT NULL REFERENCES public.repository_tb(id),
    url TEXT NOT NULL,
    expanded_url TEXT,
    file TEXT NOT NULL,
    line_number INTEGER NOT NULL,
    github_url TEXT,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (repo_id, url, file, line_number)
);
