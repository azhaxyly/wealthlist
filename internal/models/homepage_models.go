package models

type HomePageDto struct {
	TopMillionaires []Millionaire `json:"topMillionaires"`
	Featured        []Millionaire `json:"featured,omitempty"`
}
