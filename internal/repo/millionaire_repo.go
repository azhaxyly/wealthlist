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
	ScanRows(rows *sql.Rows) ([]models.Millionaire, error)
}

type MillionaireFilter struct {
	LastName   string
	FirstName  string
	MiddleName string
	Country    string
}

type millionaireRepo struct {
	db *sql.DB
}

const (
	baseQuery  = `SELECT id, last_name, first_name, middle_name, birth_date, birth_place, company, net_worth, industry, country, path_to_photo, created_at, updated_at FROM millionaires`
	countQuery = `SELECT COUNT(*) FROM millionaires`
)

func NewMillionaireRepo(db *sql.DB) *millionaireRepo {
	return &millionaireRepo{db: db}
}

func (r *millionaireRepo) Create(m *models.Millionaire) error {
	query := `
    INSERT INTO millionaires (
        last_name, first_name, middle_name, birth_date, 
        birth_place, company, net_worth, industry, 
        country, path_to_photo, created_at, updated_at
    ) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) 
    RETURNING id`

	err := r.db.QueryRow(query,
		m.LastName, m.FirstName, m.MiddleName, m.BirthDate,
		m.BirthPlace, m.Company, m.NetWorth, m.Industry,
		m.Country, m.PathToPhoto,
	).Scan(&m.ID)

	return err
}

func (r *millionaireRepo) GetByID(id int) (*models.Millionaire, error) {
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

func (r *millionaireRepo) Update(m *models.Millionaire) error {
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

func (r *millionaireRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM millionaires WHERE id = $1", id)
	return err
}

func (r *millionaireRepo) Search(filter MillionaireFilter, page int, pageSize int) (models.PaginationMillionaireDto, error) {
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

func (r *millionaireRepo) GetAll(page int, pageSize int) (models.PaginationMillionaireDto, error) {
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
