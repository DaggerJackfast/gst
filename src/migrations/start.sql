-- CREATE DATABASE gst;

CREATE TABLE users
(
    id       integer primary key generated by default as identity,
    email    varchar(255) not null unique,
    username varchar(255) not null unique,
    password varchar(255) not null
);

CREATE TABLE user_profile_tokens
(
    id            integer primary key generated by default as identity,
    user_id       integer              not null,
    profile_token varchar(255)         not null,
    token_type    varchar(50)          not null,
    is_active     boolean default true not null,
    expired_in    integer default 0    not null,
    created_at    timestamp            not null,
    updated_at    timestamp            not null,
    foreign key (user_id) references users (id)
);