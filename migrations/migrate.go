package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrationUp(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = 'millionaires'
		);
	`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}

	if exists {
		log.Println("Table 'millionaires' already exists, skipping creation")
		return nil
	}

	upSQL := `
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
	`
	_, err = db.Exec(upSQL)
	return err
}

func RunMigrationDown(db *sql.DB) error {
	downSQL := `
	DROP TABLE IF EXISTS millionaires;
	`
	_, err := db.Exec(downSQL)
	return err
}
