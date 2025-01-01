CREATE TABLE Categories(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL DEFAULT(''),
    icon text NOT NULL,
    slug text NOT NULL
);