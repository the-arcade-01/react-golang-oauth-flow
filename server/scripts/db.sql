create database if not exists auth_flow;

use auth_flow;

create table if not exists users (
    user_id serial PRIMARY KEY,
    google_id varchar(255) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    name varchar(255),
    picture varchar(255)
);