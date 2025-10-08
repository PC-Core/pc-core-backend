CREATE TABLE IF NOT EXISTS MouseChars(
  id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name text NOT NULL,
  type text NOT NULL CHECK (type IN ('мышь', 'тачпад')),
  dpi integer NOT NULL,
  release_year integer NOT NULL
)

ALTER TABLE MouseChars ADD COLUMN mouse_id integer REFERENCES MouseChars(id);

INSERT INTO Categories (title, description, icon, slug)
VALUES ('Мышь', 'Компьютерные мыши', 'mouse-icon', 'mouse');