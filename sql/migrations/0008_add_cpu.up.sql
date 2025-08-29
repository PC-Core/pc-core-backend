CREATE TABLE IF NOT EXISTS CpuChars(
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    pcores integer NOT NULL,
    ecores integer DEFAULT 0 NOT NULL,
    threads integer NOT NULL,
    base_p_freq_mhz integer DEFAULT 0 NOT NULL,
    max_p_freq_mhz integer DEFAULT 0 NOT NULL,
    base_e_freq_mhz integer DEFAULT 0 NOT NULL,
    max_e_freq_mhz integer DEFAULT 0 NOT NULL,
    socket text NOT NULL,
    l1_kb integer NOT NULL,
    l2_kb integer NOT NULL,
    l3_kb integer NOT NULL,
    tecproc_nm integer NOT NULL,
    tdp_watt integer NOT NULL,
    release_year integer NOT NULL
);

ALTER TABLE LaptopChars ALTER COLUMN cpu SET DATA TYPE integer USING cpu::integer;
ALTER TABLE LaptopChars RENAME COLUMN cpu TO cpu_id;
