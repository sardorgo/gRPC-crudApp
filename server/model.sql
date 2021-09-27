-- The name of db is grpcapp

create table users (
    user_id uuid not null primary key,
    first_name varchar(32) not null,
    last_name varchar(32) not null
);