CREATE TABLE `order`
(
    id BINARY(16) NOT NULL,
    cost FLOAT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE order_item
(
    order_id BINARY(16) NOT NULL,
    menu_item_id BINARY(16) NOT NULL,
    quantity INT,
    PRIMARY KEY (order_id, menu_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
