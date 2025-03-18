CREATE TYPE CpuSocket AS ENUM (
    'AM4', 'AM5', 
    'LGA775', 'LGA1156', 'LGA1155', 'LGA1150', 
    'LGA1151', 'LGA1151v2', 'LGA1200', 'LGA1700', 'LGA1851'
);

CREATE TABLE IF NOT EXISTS CpuChars(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    pcores integer NOT NULL,
    ecores integer DEFAULT 0 NOT NULL,
    threads integer NOT NULL,
    base_freq_mhz integer DEFAULT 0 NOT NULL,
    max_freq_mhz integer DEFAULT 0 NOT NULL,
    socket CpuSocket NOT NULL,
    l1_kb integer NOT NULL,
    l2_kb integer NOT NULL,
    l3_kb integer NOT NULL,
    tecproc_nm integer NOT NULL,
    tdp_watt integer NOT NULL,
    release_year integer NOT NULL
);

ALTER TABLE LaptopChars ALTER COLUMN cpu integer;
ALTER TABLE LaptopChars RENAME COLUMN cpu cpu_id;