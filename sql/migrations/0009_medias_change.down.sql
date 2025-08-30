ALTER TABLE products ADD COLUMN medias integer[];

UPDATE products SET medias = (
    SELECT json_agg(id) FROM medias WHERE medias.product_id = products.id
);

ALTER TABLE medias DROP CONSTRAINT fk_product;
ALTER TABLE medias DROP COLUMN product_id;