ALTER TABLE shopping_cart_items
DROP COLUMN IF EXISTS final_price;

ALTER TABLE shopping_cart_items
DROP COLUMN IF EXISTS discount;
