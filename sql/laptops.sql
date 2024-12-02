CREATE TABLE laptops (
	id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name text NOT NULL,
	cpu text NOT NULL,
	ram smallint NOT NULL,
	gpu text NOT NULL,
	price numeric(18, 2) NOT NULL,
	discount smallint NOT NULL DEFAULT 0 CHECK (Discount >= 0 AND Discount <= 100)
);