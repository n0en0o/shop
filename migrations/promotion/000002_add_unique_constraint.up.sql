DELETE p1
FROM promos p1
    INNER JOIN promos p2
WHERE
    p1.catalog_item_id = p2.catalog_item_id
    AND p1.id > p2.id;

ALTER TABLE promos
ADD CONSTRAINT unique_catalog_item_id UNIQUE (catalog_item_id);
