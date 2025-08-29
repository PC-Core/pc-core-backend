CREATE TYPE MediaType AS ENUM ('Image', 'Video');

CREATE TABLE IF NOT EXISTS Medias(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    url text not null,
    type MediaType
);

ALTER TABLE Products ADD Medias integer[];