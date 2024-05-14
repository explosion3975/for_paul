CREATE TABLE customer (
    id CHAR(20) NOT NULL,
    password CHAR(64) NOT NULL,
    name VARCHAR(15) NOT NULL,
    phone CHAR(10) NOT NULL,
    address VARCHAR(256) NOT NULL,
    age INT NOT NULL,
    job CHAR(20) NOT NULL,
    join_date DATE NOT NULL DEFAULT (CURRENT_DATE),
    image CHAR(50) NOT NULL,
    premission INT NOT NULL DEFAULT 0,
    status BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);

CREATE TABLE product_record(
    product_id INT NOT NULL AUTO_INCREMENT,
    seller_id CHAR(10) NOT NULL,
    product_name CHAR(128),
    image CHAR(50) NOT NULL,
    item_sold INT NOT NULL DEFAULT 0,
    price INT NOT NULL,
    number INT NOT NULL,
    location CHAR(32),
    PRIMARY KEY (product_id),
    FOREIGN KEY (seller_id) REFERENCES customer(id)
);

CREATE TABLE order_record(
    id CHAR(20) NOT NULL,
    order_id INT NOT NULL AUTO_INCREMENT,
    product_id INT NOT NULL,
    order_date DATE NOT NULL,
    number INT NOT NULL,
    seller_id CHAR(10) NOT NULL,
    PRIMARY KEY (order_id),
    FOREIGN KEY (seller_id) REFERENCES customer(id),
    FOREIGN KEY (product_id) REFERENCES product_record(product_id)
    
);
