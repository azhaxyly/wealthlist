package repo

import (
	"database/sql"
	"millionaire-list/internal/models"
)

func (r *millionaireRepo) ScanRows(rows *sql.Rows) ([]models.Millionaire, error) {
	var millionaires []models.Millionaire
	for rows.Next() {
		var m models.Millionaire
		err := rows.Scan(
			&m.ID, &m.LastName, &m.FirstName, &m.MiddleName,
			&m.BirthDate, &m.BirthPlace, &m.Company, &m.NetWorth,
			&m.Industry, &m.Country, &m.PathToPhoto,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		millionaires = append(millionaires, m)
	}
	return millionaires, rows.Err()
}
