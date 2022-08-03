create table users
(
    id varchar(255) not null primary key,
    username varchar(255) not null unique comment 'username'
);

