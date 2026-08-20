-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    account_name VARCHAR(100) NOT NULL,
    total_amount DECIMAL(18,2) NOT NULL DEFAULT 0,

    -- Order Status
    current_order_status VARCHAR(50) NOT NULL DEFAULT 'Draft',

    -- Contact Info (embedded)
    contact_first_name VARCHAR(100) NOT NULL,
    contact_last_name VARCHAR(100) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,

    -- Delivery Address (embedded)
    address_street VARCHAR(200) NOT NULL,
    address_city VARCHAR(100) NOT NULL,
    address_region VARCHAR(100) NOT NULL,
    address_postal_code VARCHAR(20) NOT NULL,

    -- Payment
    current_payment_method VARCHAR(50) NOT NULL DEFAULT 'CreditCard',
    current_payment_status VARCHAR(50) NOT NULL DEFAULT 'Pending',

    -- Card Details (optional, embedded)
    card_name VARCHAR(100),
    card_number VARCHAR(20),
    card_expiration VARCHAR(10),
    card_cvv VARCHAR(10),

    -- Audit
    created_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_modified_by VARCHAR(100),
    last_modified_at TIMESTAMP WITH TIME ZONE
);

-- Order Items table
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    catalog_item_name VARCHAR(200) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(18, 2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_account_name ON orders (account_name);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (current_order_status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders (current_payment_status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

-- Seed data
INSERT INTO orders (
    id,
    account_name,
    total_amount,
    current_order_status,
    contact_first_name,
    contact_last_name,
    contact_email,
    address_street,
    address_city,
    address_region,
    address_postal_code,
    current_payment_method,
    current_payment_status,
    card_name,
    card_number,
    card_expiration,
    card_cvv,
    created_by,
    last_modified_by
)
VALUES
    (
        '28282828-2828-2828-2828-282828282828',
        'demo-workstation',
        1599.98,
        'Submitted',
        'Alex',
        'Morgan',
        'alex.morgan@example.com',
        '12 Tech Park Lane',
        'Seattle',
        'WA',
        '98101',
        'CreditCard',
        'Paid',
        'Alex Morgan',
        '4111111111111111',
        '12/30',
        '123',
        'seed',
        'seed'
    ),
    (
        '29292929-2929-2929-2929-292929292929',
        'demo-media',
        1649.98,
        'Processing',
        'Jamie',
        'Taylor',
        'jamie.taylor@example.com',
        '45 Media Center Drive',
        'Austin',
        'TX',
        '78701',
        'CreditCard',
        'Paid',
        'Jamie Taylor',
        '5555555555554444',
        '08/29',
        '456',
        'seed',
        'seed'
    );

INSERT INTO order_items (
    order_id,
    catalog_item_name,
    quantity,
    unit_price
)
VALUES
    (
        '28282828-2828-2828-2828-282828282828',
        'Lenovo ThinkPad X1',
        1,
        1499.99
    ),
    (
        '28282828-2828-2828-2828-282828282828',
        'Logitech MX Master 3S',
        1,
        99.99
    ),
    (
        '29292929-2929-2929-2929-292929292929',
        'LG OLED C3',
        1,
        1299.99
    ),
    (
        '29292929-2929-2929-2929-292929292929',
        'Sony WH-1000XM5',
        1,
        349.99
    );
