CREATE TABLE IF NOT EXISTS Products(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    price numeric NOT NULL,
    selled bigint NOT NULL DEFAULT(0) CHECK(selled >= 0),
    stock bigint NOT NULL DEFAULT(0) CHECK(stock >= 0),
    chars_table_name text NOT NULL,
    chars_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS LaptopChars(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cpu text NOT NULL,
    ram smallint NOT NULL,
    gpu text NOT NULL
);

CREATE TYPE UserRole AS ENUM ('Default', 'Admin');

CREATE TABLE Users (
	id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	Name varchar(30) UNIQUE NOT NULL,
	Email text UNIQUE NOT NULL CHECK (Email ~* '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
	Role UserRole NOT NULL DEFAULT 'Default',
	PasswordHash text NOT NULL
);