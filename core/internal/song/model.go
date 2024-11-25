package song

import (
	"database/sql"
	"time"
)

type SongModel struct {
	Id          string         `db:"id"`
	Song        string         `db:"song"`
	Group       string         `db:"group"`
	Link        sql.NullString `db:"link"`
	Text        sql.NullString `db:"text"`
	ReleaseDate sql.NullTime   `db:"release_date"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type SongShowModel struct {
	Id          string     `json:"id"`
	Song        string     `json:"song"`
	Group       string     `json:"group"`
	Link        *string    `json:"link"`
	Text        *string    `json:"text"`
	ReleaseDate *time.Time `json:"releaseDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (m *SongModel) MapToShow() SongShowModel {

	var link *string = nil
	var text *string = nil
	var releaseDate *time.Time = nil

	if m.Link.Valid {
		link = &m.Link.String
	}

	if m.Text.Valid {
		text = &m.Text.String
	}

	if m.ReleaseDate.Valid {
		releaseDate = &m.ReleaseDate.Time

	}

	return SongShowModel{
		Id:          m.Id,
		Song:        m.Song,
		Group:       m.Group,
		Link:        link,
		Text:        text,
		ReleaseDate: releaseDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func NewSong(song, group, link string, releaseDate time.Time) SongModel {
	return SongModel{
		Song:        song,
		Group:       group,
		Link:        sql.NullString{String: link},
		ReleaseDate: sql.NullTime{Time: releaseDate},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func MapModelsToShow(models []SongModel) []SongShowModel {
	songsShow := []SongShowModel{}
	for _, v := range models {
		songsShow = append(songsShow, v.MapToShow())
	}
	return songsShow
}
