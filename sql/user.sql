CREATE TABLE users (
	id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	Name varchar(30) UNIQUE NOT NULL,
	Email text UNIQUE NOT NULL CHECK (Email ~* '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
	Role UserRole NOT NULL DEFAULT 'Default',
	PasswordHash text NOT NULL
);