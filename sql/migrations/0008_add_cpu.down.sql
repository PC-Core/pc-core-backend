DROP TABLE CpuChars;

ALTER TABLE LaptopChars RENAME COLUMN cpu_id cpu;
ALTER TABLE LaptopChars ALTER COLUMN cpu text;