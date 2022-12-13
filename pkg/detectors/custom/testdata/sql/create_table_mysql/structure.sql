-- MySQL CREATE TABLE statement
CREATE TABLE users (
       id integer NOT NULL,
       class_id integer,
       type char(255) NOT NULL,
       created_at timestamp NOT NULL,
       updated_at timestamp NOT NULL
);
