ALTER TABLE LaptopChars DROP COLUMN gpu_id;

DROP TABLE IF EXISTS GpuChars;

DELETE FROM Categories WHERE slug = 'gpu';