package couplet

import "time"

type CoupletModel struct {
	Id         string    `db:"id" json:"id"`
	SongId     string    `db:"song_id" json:"-"`
	CoupletNum int       `db:"couplet_num" json:"coupletNum"`
	Couplet    string    `db:"couplet" json:"couplet"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
}

func NewCouplet(songId, couplet string, coupletNum int) *CoupletModel {
	return &CoupletModel{
		SongId:     songId,
		CoupletNum: coupletNum,
		Couplet:    couplet,
		CreatedAt:  time.Now(),
	}
}
