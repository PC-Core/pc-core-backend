CREATE TABLE IF NOT EXISTS GpuChars(
  id integer GENERATE ALWAYS AS IDENTITY PRIMARY KEY,
  name text NOT NULL,
  memory_gb integer NOT NULL,
  memory_type text NOT NULL,
  bus_width_bit integer NOT NULL,
  base_freq_mhz integer NOT NULL,
  boost_freq_mhz integer NOT NULL,
  tecproc_nm integer NOT NULL,
  tdp_watt integer NOT NULL,
  release_year integer NOT NULL
);

ALTER TABLE LaptopChars ADD COLUMN gpu_id integer REFERENCES GpuChars(id);
INSERT INTO Categories (title, description, icon, slug) 
VALUES ('GPU', 'Графические процессоры', 'gpu-icon', 'gpu', 'Видеокарта')
ON CONFLICT (slug) DO NOTHING;