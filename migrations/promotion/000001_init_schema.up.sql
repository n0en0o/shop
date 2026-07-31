CREATE TABLE IF NOT EXISTS promos (
    id CHAR(36) PRIMARY KEY,
    catalog_item_id VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    value DECIMAL(10, 2) NOT NULL DEFAULT 0
);

-- Посевочные данные
INSERT INTO promos (
    id,
    catalog_item_id,
    title,
    value
)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'a0000001-0000-0000-0000-000000000001',
        'Скидка 1000 денег',
        1000.00
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'a0000001-0000-0000-0000-000000000002',
        'Скидка 2000 денег',
        2000.00
    ),
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'a0000001-0000-0000-0000-000000000004',
        'Скидка 100 денег',
        100.00
    );
