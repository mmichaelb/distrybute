-- users ddl
CREATE TABLE IF NOT EXISTS distrybute.users (
    id              uuid,
    username        varchar(32)     NOT NULL,
    auth_token      text            NULL,
    password_alg    varchar(32)     NOT NULL,
    password_salt   bytea           NOT NULL,
    "password"      bytea           NOT NULL,
    CONSTRAINT users_pk                 PRIMARY KEY (id),
    CONSTRAINT users_auth_token_unique  UNIQUE (auth_token)
);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_un ON distrybute.users (LOWER(username));

-- files db ddl
CREATE TABLE IF NOT EXISTS distrybute.entries (
    id                  uuid,
    author              uuid            NOT NULL,
    call_reference      varchar(4)      NOT NULL,
    delete_reference    varchar(12)     NOT NULL,
    filename            varchar(256)    NOT NULL,
    content_type        varchar(127)    NOT NULL,
    upload_date         timestamptz     NOT NULL,
    "size"              bigint,
    CONSTRAINT entries_pk                       PRIMARY KEY (id),
    CONSTRAINT entries_fk                       FOREIGN KEY (author) REFERENCES distrybute.users(id),
    CONSTRAINT entries_call_reference_unique    UNIQUE (call_reference),
    CONSTRAINT entries_delete_reference_unique  UNIQUE (delete_reference)
);
