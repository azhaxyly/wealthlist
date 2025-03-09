package repo

import (
	"database/sql"
	"fmt"
	"millionaire-list/internal/models"
)

type MillionaireRepository interface {
	Create(m *models.Millionaire) error
	GetByID(id int) (*models.Millionaire, error)
	Search(filter MillionaireFilter, page int, pageSize int) (models.PaginationMillionaireDto, error)
	GetAll(page int, pageSize int) (models.PaginationMillionaireDto, error)
	Update(m *models.Millionaire) error
	Delete(id int) error
	GetWithPhotos() ([]models.Millionaire, error)
	UpdatePhotoPath(id int, path string) error
	BatchUpdatePhotoPaths(photoMap map[int]string) error
}

type MillionaireFilter struct {
	LastName   string
	FirstName  string
	MiddleName string
	Country    string
}

type MillionaireRepo struct {
	db *sql.DB
}

const (
	baseQuery  = `SELECT id, last_name, first_name, middle_name, birth_date, birth_place, company, net_worth, industry, country, path_to_photo, created_at, updated_at FROM millionaires`
	countQuery = `SELECT COUNT(*) FROM millionaires`
)

func NewMillionaireRepo(db *sql.DB) *MillionaireRepo {
	return &MillionaireRepo{db: db}
}

func (r *MillionaireRepo) Create(m *models.Millionaire) error {
	query := `
    INSERT INTO millionaires (
        last_name, first_name, middle_name, birth_date, 
        birth_place, company, net_worth, industry, 
        country, path_to_photo, created_at, updated_at
    ) 
    VALUES ($1, $2, $3, NULLIF($4, '')::DATE, $5, $6, $7, $8, $9, $10, NOW(), NOW()) 
    RETURNING id`

	err := r.db.QueryRow(query,
		m.LastName, m.FirstName, m.MiddleName, m.BirthDate,
		m.BirthPlace, m.Company, m.NetWorth, m.Industry,
		m.Country, m.PathToPhoto,
	).Scan(&m.ID)

	return err
}

func (r *MillionaireRepo) GetByID(id int) (*models.Millionaire, error) {
	query := baseQuery + " WHERE id = $1"
	row := r.db.QueryRow(query, id)

	m := &models.Millionaire{}
	err := row.Scan(
		&m.ID, &m.LastName, &m.FirstName, &m.MiddleName,
		&m.BirthDate, &m.BirthPlace, &m.Company, &m.NetWorth,
		&m.Industry, &m.Country, &m.PathToPhoto,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MillionaireRepo) Update(m *models.Millionaire) error {
	query := `
		UPDATE millionaires 
		SET last_name = $1, first_name = $2, middle_name = $3,
		    birth_date = $4, birth_place = $5, company = $6,
		    net_worth = $7, industry = $8, country = $9,
		    path_to_photo = $10, updated_at = NOW()
		WHERE id = $11`

	_, err := r.db.Exec(query,
		m.LastName, m.FirstName, m.MiddleName, m.BirthDate,
		m.BirthPlace, m.Company, m.NetWorth, m.Industry,
		m.Country, m.PathToPhoto, m.ID,
	)
	return err
}

func (r *MillionaireRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM millionaires WHERE id = $1", id)
	return err
}

func (r *MillionaireRepo) GetWithPhotos() ([]models.Millionaire, error) {
	query := baseQuery + " WHERE path_to_photo IS NOT NULL AND path_to_photo != ''"
	return r.FetchMillionaires(query)
}

func (r *MillionaireRepo) UpdatePhotoPath(id int, path string) error {
	_, err := r.db.Exec(
		"UPDATE millionaires SET path_to_photo = $1 WHERE id = $2",
		path, id,
	)
	return err
}

func (r *MillionaireRepo) BatchUpdatePhotoPaths(photoMap map[int]string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		UPDATE millionaires 
		SET path_to_photo = $1 
		WHERE id = $2`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for id, path := range photoMap {
		if _, err := stmt.Exec(path, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *MillionaireRepo) Search(filter MillionaireFilter, page int, pageSize int) (models.PaginationMillionaireDto, error) {
	result := models.PaginationMillionaireDto{
		Page:     page,
		PageSize: pageSize,
	}

	where, args := BuildWhereClause(filter)
	query := baseQuery + where + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return result, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	millionaires, err := r.ScanRows(rows)
	if err != nil {
		return result, err
	}

	total, err := r.GetTotalCount(where, args...)
	if err != nil {
		return result, err
	}

	result.Millionaires = millionaires
	result.Total = total
	return result, nil
}

func (r *MillionaireRepo) GetAll(page int, pageSize int) (models.PaginationMillionaireDto, error) {
	result := models.PaginationMillionaireDto{
		Page:     page,
		PageSize: pageSize,
	}

	query := baseQuery + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(query)
	if err != nil {
		return result, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	millionaires, err := r.ScanRows(rows)
	if err != nil {
		return result, err
	}

	total, err := r.GetTotalCount("")
	if err != nil {
		return result, err
	}

	result.Millionaires = millionaires
	result.Total = total
	return result, nil
}
