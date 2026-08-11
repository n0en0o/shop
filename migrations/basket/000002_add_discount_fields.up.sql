ALTER TABLE shopping_cart_items
ADD COLUMN IF NOT EXISTS discount NUMERIC(10,2) NOT NULL DEFAULT 0;

ALTER TABLE shopping_cart_items
ADD COLUMN IF NOT EXISTS final_price NUMERIC(10,2) NOT NULL DEFAULT 0;

UPDATE shopping_cart_items
SET final_price = unit_price
WHERE final_price = 0;
