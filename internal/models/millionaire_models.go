package models

import "time"

type Millionaire struct {
	ID          int        `json:"id"`
	LastName    string     `json:"lastName"`
	FirstName   string     `json:"firstName"`
	MiddleName  *string    `json:"middleName,omitempty"` // Указатель, чтобы пустое поле не сериализовалось
	BirthDate   *time.Time `json:"birthDate,omitempty"`
	BirthPlace  *string    `json:"birthPlace,omitempty"`
	Company     *string    `json:"company,omitempty"`
	NetWorth    float64    `json:"netWorth,omitempty"`
	Industry    *string    `json:"industry,omitempty"`
	Country     *string    `json:"country,omitempty"`
	PathToPhoto *string    `json:"pathToPhoto,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}
