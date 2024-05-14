CREATE TABLE IF NOT EXISTS stocks (
    stock_id integer primary key generated always as identity,
    name varchar NOT NULL,
    is_available boolean NOT NULL DEFAULT TRUE
    );

CREATE TABLE IF NOT EXISTS products (
    product_id integer primary key generated always as identity,
    name varchar NOT NULL,
    size varchar NOT NULL
    );

CREATE TABLE IF NOT EXISTS product_on_stock (
    stock_id integer,
    product_id integer,
    available_qty bigint,
    reserved_qty bigint DEFAULT 0,
    CONSTRAINT stock_fkey
    FOREIGN KEY (stock_id) REFERENCES stocks (stock_id),
    CONSTRAINT product_fkey
    FOREIGN KEY (product_id) REFERENCES products (product_id),
    PRIMARY KEY (stock_id, product_id)
    );

INSERT INTO products
    (name, size)
VALUES
    ( 'dress', 'S'),
    ( 'jacket', 'S'),
    ( 'shoes', '38'),
    ( 'socks', '36-38'),
    ( 'trousers', 'XS');

INSERT INTO stocks
    (name, is_available)
VALUES
    ( 'scheremetevo', true ),
    ( 'domodedovo', true ),
    ( 'vnukovo', false );

INSERT INTO product_on_stock
    (stock_id, product_id, available_qty)
VALUES
    ( 1, 1, 5),
    ( 1, 2, 3 ),
    ( 1, 3, 10 ),
    ( 2, 4, 1 );
    