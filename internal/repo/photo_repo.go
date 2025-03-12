package repo

import (
	"database/sql"
	"log/slog"
	"wealthlist/internal/logger"
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

func (r *PhotoRepo) GetPhotoPath(millionaireID int) (string, error) {
	var photoPath string
	query := `SELECT path_to_photo FROM millionaires WHERE id = $1`
	err := r.DB.QueryRow(query, millionaireID).Scan(&photoPath)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		r.log.Error("Error fetching photo path", logger.Err(err))
		return "", err
	}
	return photoPath, nil
}

func (r *PhotoRepo) ClearPhotoPath(millionaireID int) error {
	query := `UPDATE millionaires SET path_to_photo = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.DB.Exec(query, millionaireID)
	if err != nil {
		r.log.Error("Error clearing photo path", logger.Err(err))
		return err
	}
	return nil
}
