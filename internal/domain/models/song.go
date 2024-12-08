package models

import (
	"github.com/google/uuid"
	"github.com/hardfinhq/go-date"
)

type Song struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate date.Date `json:"release_date" swaggertype:"primitive,string"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}
