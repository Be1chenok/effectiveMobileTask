CREATE TABLE IF NOT EXISTS persons(
id SERIAL PRIMARY KEY,
name VARCHAR(64) NOT NULL,
surname VARCHAR(64) NOT NULL,
patronymic VARCHAR(64),
age SMALLINT NOT NULL,
gender VARCHAR(6) NOT NULL,
nationality VARCHAR(5) NOT NULL
);