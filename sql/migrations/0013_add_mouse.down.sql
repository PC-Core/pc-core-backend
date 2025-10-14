ALTER TABLE MouseChars DROP COLUMN mouse_id;

DROP TABLE IF EXISTS MouseChars;

DELETE FROM Categories WHERE slug = 'mouse';