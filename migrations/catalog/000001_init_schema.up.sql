CREATE TABLE brands (
    id UUID PRIMARY KEY,
    title VARCHAR(80) NOT NULL
);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    title VARCHAR(80) NOT NULL
);

CREATE TABLE catalog_items (
    id UUID PRIMARY KEY,
    title VARCHAR(80) NOT NULL,
    short_description TEXT,
    full_description TEXT,
    image_url TEXT,
    brand_id UUID REFERENCES brands (id),
    category_id UUID REFERENCES categories (id),
    price DOUBLE PRECISION NOT NULL DEFAULT 0
);


INSERT INTO brands (id, title) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Apple'),
    ('22222222-2222-2222-2222-222222222222', 'Samsung'),
    ('33333333-3333-3333-3333-333333333333', 'Sony'),
    ('44444444-4444-4444-4444-444444444444', 'LG'),
    ('55555555-5555-5555-5555-555555555555', 'Lenovo'),
    ('66666666-6666-6666-6666-666666666666', 'Dell'),
    ('77777777-7777-7777-7777-777777777777', 'HP'),
    ('88888888-8888-8888-8888-888888888888', 'Asus'),
    ('99999999-9999-9999-9999-999999999999', 'Xiaomi'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Logitech');

INSERT INTO categories (id, title) VALUES
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Smartphones'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Laptops'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Tablets'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Headphones'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'TVs'),
    ('12121212-1212-1212-1212-121212121212', 'Monitors'),
    ('13131313-1313-1313-1313-131313131313', 'Keyboards'),
    ('14141414-1414-1414-1414-141414141414', 'Mice'),
    ('15151515-1515-1515-1515-151515151515', 'Cameras'),
    ('16161616-1616-1616-1616-161616161616', 'Speakers');

INSERT INTO catalog_items (
    id, 
    title, 
    short_description, 
    full_description, 
    image_url, 
    brand_id, 
    category_id, 
    price) 
VALUES
    ('17171717-1717-1717-1717-171717171717', 'iPhone 15', 'Apple smartphone', 'Apple iPhone 15 with OLED display and fast performance.', 'https://example.com/images/iphone-15.jpg', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 799.99),
    ('18181818-1818-1818-1818-181818181818', 'Galaxy S24', 'Samsung smartphone', 'Samsung Galaxy S24 with high refresh rate display.', 'https://example.com/images/galaxy-s24.jpg', '22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 749.99),
    ('19191919-1919-1919-1919-191919191919', 'Sony WH-1000XM5', 'Noise cancelling headphones', 'Wireless Sony headphones with active noise cancellation.', 'https://example.com/images/sony-wh-1000xm5.jpg', '33333333-3333-3333-3333-333333333333', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 349.99),
    ('20202020-2020-2020-2020-202020202020', 'LG OLED C3', 'OLED TV', 'LG OLED C3 television with vivid contrast and smart TV features.', 'https://example.com/images/lg-oled-c3.jpg', '44444444-4444-4444-4444-444444444444', 'ffffffff-ffff-ffff-ffff-ffffffffffff', 1299.99),
    ('21212121-2121-2121-2121-212121212121', 'Lenovo ThinkPad X1', 'Business laptop', 'Lightweight Lenovo ThinkPad laptop for business productivity.', 'https://example.com/images/thinkpad-x1.jpg', '55555555-5555-5555-5555-555555555555', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 1499.99),
    ('23232323-2323-2323-2323-232323232323', 'Dell XPS 13', 'Compact laptop', 'Dell XPS 13 ultrabook with compact design and bright display.', 'https://example.com/images/dell-xps-13.jpg', '66666666-6666-6666-6666-666666666666', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 1199.99),
    ('24242424-2424-2424-2424-242424242424', 'HP Spectre x360', 'Convertible laptop', 'HP Spectre x360 convertible laptop with touch display.', 'https://example.com/images/hp-spectre-x360.jpg', '77777777-7777-7777-7777-777777777777', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 1099.99),
    ('25252525-2525-2525-2525-252525252525', 'Asus ROG Swift', 'Gaming monitor', 'Asus ROG Swift monitor for smooth gaming and sharp visuals.', 'https://example.com/images/asus-rog-swift.jpg', '88888888-8888-8888-8888-888888888888', '12121212-1212-1212-1212-121212121212', 599.99),
    ('26262626-2626-2626-2626-262626262626', 'Xiaomi Pad 6', 'Android tablet', 'Xiaomi Pad 6 tablet for entertainment and everyday work.', 'https://example.com/images/xiaomi-pad-6.jpg', '99999999-9999-9999-9999-999999999999', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 329.99),
    ('27272727-2727-2727-2727-272727272727', 'Logitech MX Master 3S', 'Wireless mouse', 'Logitech MX Master 3S ergonomic wireless mouse for productivity.', 'https://example.com/images/mx-master-3s.jpg', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '14141414-1414-1414-1414-141414141414', 99.99);
