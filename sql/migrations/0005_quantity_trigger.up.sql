CREATE OR REPLACE FUNCTION check_quantity_stock()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity > (SELECT stock FROM Product WHERE product_id = NEW.product_id) THEN
        RAISE EXCEPTION 'Quantity exceeds available stock';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_cart_quantity
BEFORE INSERT OR UPDATE ON Cart
FOR EACH ROW
EXECUTE FUNCTION check_quantity_stock();
