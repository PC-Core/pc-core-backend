CREATE TABLE laptops (
	id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name text NOT NULL,
	cpu text NOT NULL,
	ram smallint NOT NULL,
	gpu text NOT NULL,
	price money NOT NULL
);