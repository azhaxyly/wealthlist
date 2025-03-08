CREATE TABLE millionaires (
    id SERIAL PRIMARY KEY,
    last_name VARCHAR(500) NOT NULL,
    first_name VARCHAR(500) NOT NULL,
    middle_name VARCHAR(500),
    birth_date DATE,
    birth_place TEXT,
    net_worth BIGINT NOT NULL, -- Чистый капитал в долларах
    industry TEXT, -- Индустрия, в которой заработал состояние
    country TEXT, -- Страна проживания
    company TEXT, -- Основная компания
    biography TEXT, -- Краткая биография
    path_to_photo TEXT, -- Путь к файлу с фото
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);