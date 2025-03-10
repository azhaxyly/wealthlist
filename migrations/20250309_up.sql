CREATE TABLE millionaires (
    id SERIAL PRIMARY KEY,
    last_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),
    birth_date VARCHAR(10),  -- храним дату как строку
    birth_place VARCHAR(255),
    company VARCHAR(255),
    net_worth VARCHAR(50),  -- Число в виде строки
    industry VARCHAR(255),
    country VARCHAR(255),
    path_to_photo VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
