ALTER TABLE LaptopChars DROP COLUMN keyboard_id;

DROP TABLE IF EXISTS KeyboardChars;

DELETE FROM Categories WHERE slug = 'keyboard';