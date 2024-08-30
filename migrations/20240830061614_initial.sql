-- +goose Up
-- +goose StatementBegin
create table passwords
(
    id             serial
        constraint passwords_pk
            primary key,
    hash_base64    varchar not null,
    argon2_version integer not null,
    argon2_type    integer not null,
    salt_base64    varchar not null,
    time           integer not null,
    memory         integer not null,
    threads        integer not null,
    keylen         integer not null
);

create table users
(
    id          serial
        constraint users_pk
            primary key,
    name        varchar not null
        constraint users_pk2
            unique,
    password_id integer not null
        constraint users_passwords_id_fk
            references passwords
);

create table notes
(
    id          serial
        constraint notes_pk
            primary key,
    content     varchar not null,
    author_name varchar not null
        constraint notes_users_name_fk
            references users (name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table passwords;
drop table users;
drop table notes;
-- +goose StatementEnd
