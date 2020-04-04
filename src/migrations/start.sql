CREATE TABLE users
(
    id       serial primary key,
    email    varchar(255) not null unique,
    username varchar(255) not null unique,
    password varchar(255) not null
);

CREATE TABLE user_profile_tokens
(
    id            serial primary key,
    user_id       integer              not null,
    profile_token varchar(255)         not null,
    token_type    varchar(50)          not null,
    is_active     boolean default true not null,
    expired_in    integer default 0    not null,
    created_at    timestamp            not null,
    updated_at    timestamp            not null,
    foreign key (user_id) references users (id)
);

CREATE TABLE sessions
(
    id            serial primary key,
    user_id       integer      not null,
    refresh_token varchar(255) not null,
    user_agent    varchar(255) not null,
    fingerprint   varchar(255) not null,
    ip            varchar(15),
    expired_in    bigint       not null,
    created_at    timestamp    not null,
    updated_at    timestamp    not null,
    foreign key (user_id) references users (id)
);
