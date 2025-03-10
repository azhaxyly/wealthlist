package repo

import (
	"database/sql"
	"log/slog"
	"millionaire-list/internal/logger"
)

type PhotoRepo struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewPhotoRepo(db *sql.DB, log *slog.Logger) *PhotoRepo {
	return &PhotoRepo{
		DB:  db,
		log: log,
	}
}

// Обновляет путь к фото у миллионера
func (r *PhotoRepo) UpdatePhotoPath(millionaireID int, filePath string) error {
	query := `
		UPDATE millionaires
		SET path_to_photo = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.DB.Exec(query, filePath, millionaireID)
	if err != nil {
		r.log.Error("Error updating millionaire photo path", logger.Err(err))
		return err
	}
	return nil
}
