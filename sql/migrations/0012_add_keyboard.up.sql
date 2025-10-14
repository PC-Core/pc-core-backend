CREATE TABLE IF NOT EXISTS KeyboardChars(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    type text NOT NULL CHECK (type IN ('механическая', 'мембранная')),
    switches text[] NOT NULL,
    release_year integer NOT NULL
);

ALTER TABLE LaptopChars ADD COLUMN keyboard_id integer REFERENCES KeyboardChars(id);

INSERT INTO Categories (title, description, icon, slug) 
VALUES ('Клавиатуры', 'Компьютерные клавиатуры', 'keyboard-icon', 'keyboard');