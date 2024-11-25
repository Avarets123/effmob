package song

import (
	"effect-mobile/pkg/postgres"
	"effect-mobile/pkg/res"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	logger *slog.Logger
	db     *postgres.PostgresDb
}

func NewRepository(logger *slog.Logger, db *postgres.PostgresDb) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}

func (r *Repository) Delete(id string) res.Error {
	return r.Update(&SongUpdateDto{
		Id:        id,
		DeletedAt: time.Now(),
	})

}

func (r *Repository) FindOne(where string, args ...any) (*SongModel, res.Error) {

	var song SongModel

	err := r.db.Get(&song, where, args...)

	if song.Id == "" {
		r.logger.Error(err.Error())
		return nil, res.NewError(err, 404)
	}

	return &song, nil

}

func (r *Repository) CheckExistsBySongOrId(song, id string) (*SongModel, res.Error) {

	where := ""
	arg := ""

	if song != "" {
		where = " \"song\" = $1"
		arg = song
	}

	if id != "" {
		where = " \"id\" = $1"
		arg = id
	}

	findOneWhere := fmt.Sprintf("SELECT id, song from songs where %s AND deleted_at IS NULL", where)

	fmt.Println(findOneWhere)

	return r.FindOne(findOneWhere, arg)
}

func (r *Repository) Update(song *SongUpdateDto) res.Error {

	if song.Id == "" {
		msg := "Id not passed for update song!"
		r.logger.Error(msg)
		return res.NewErrorWithMessage(msg, http.StatusBadRequest)
	}

	sql, args := postgres.GetUpdateSqlFromModel("songs", "id", song.Id, *song)

	result := r.db.QueryRowx(sql, args...)

	err := result.Err()
	if err != nil {
		return res.NewError(err, http.StatusInternalServerError)
	}

	return nil

}
func (r *Repository) FindCount(where string, args ...any) (int, res.Error) {

	var total int

	var err error
	if len(args) != 0 {
		err = r.db.Get(&total, where, args...)
	} else {
		err = r.db.Get(&total, where)
	}

	if err != nil {
		return total, res.NewError(err, http.StatusInternalServerError)
	}

	return total, nil

}

func (r *Repository) FindMany(where string, args ...any) ([]SongModel, res.Error) {

	var songs []SongModel

	err := r.db.Select(&songs, where, args...)

	if err != nil {
		return nil, res.NewError(err, 500)
	}

	return songs, nil

}

func (r *Repository) Create(song *SongModel) (string, res.Error) {

	hasSong, _ := r.CheckExistsBySongOrId(song.Song, "")

	if hasSong != nil {
		return "", res.NewErrorWithMessage("Song by passed name exists!", http.StatusBadRequest)
	}

	id := uuid.New().String()
	song.Id = id

	sql, args := postgres.GetInsertSqlFromModel("songs", *song)

	result := r.db.QueryRowx(sql, args...)
	err := result.Err()

	if err != nil {
		r.logger.Error(err.Error())
		return "", res.NewError(err, http.StatusInternalServerError)
	}

	return id, nil

}
