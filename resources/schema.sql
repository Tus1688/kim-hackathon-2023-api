CREATE TABLE IF NOT EXISTS users (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    username VARCHAR(32) NOT NULL UNIQUE,
    hashed_password BINARY(60) NOT NULL,
    is_admin BOOL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS customers(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(32) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    email VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS businesses(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(32) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT  CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(32) NOT NULL UNIQUE,
    description VARCHAR(255) NULL,
    business_refer BINARY(16) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (business_refer) REFERENCES businesses(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_images(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    product_refer BINARY(16) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_refer) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_analytics(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    product_id BINARY(16) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    country_name VARCHAR(255) NULL,
    city VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders(
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    product_refer BINARY(16) NOT NULL,
    quantity INT NOT NULL,
    commission_rate INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_refer) REFERENCES products(id) ON DELETE CASCADE
);
