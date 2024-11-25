package song

import (
	"time"
)

type SongCreateDto struct {
	Group       string `json:"group" validate:"required"`
	Song        string `json:"song" validate:"required"`
	Link        string `json:"link,omitempty" validate:"omitempty,url"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}

type SongUpdateDto struct {
	Id          string
	Group       string    `json:"group" db:"group"`
	Song        string    `json:"song" db:"song"`
	Link        string    `json:"link" db:"link" validate:"omitempty,url"`
	ReleaseDate string    `json:"releaseDate,omitempty" db:"release_date"`
	DeletedAt   time.Time `db:"deleted_at"`
}
