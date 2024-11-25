package couplet

import (
	"database/sql"
	"effect-mobile/pkg/postgres"
	"effect-mobile/pkg/res"
	"effect-mobile/pkg/utils"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Repository struct {
	db     *postgres.PostgresDb
	logger *slog.Logger
}

func NewRepository(db *postgres.PostgresDb, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) queryAction(trx *sql.Tx, where string, args ...any) res.Error {

	row, err := r.getExecutor(trx)(where, args...)
	row.Close()

	if err != nil {
		r.logger.Error(err.Error())
		return res.NewError(err, http.StatusInternalServerError)
	}
	return nil
}

func (r *Repository) Create(trx *sql.Tx, couplets []CoupletModel) res.Error {
	query, args := postgres.GetInsertBatchSqlFromModels("couplets", couplets)
	return r.queryAction(trx, query, args...)
}

func (r *Repository) Delete(trx *sql.Tx, songId string, coupletIds []string) res.Error {
	where := "DELETE FROM couplets WHERE song_id = $1"
	if len(coupletIds) == 0 {
		return r.queryAction(trx, where, songId)
	}

	argsStr := utils.ReduceSlice(coupletIds, func(acc string, el string, i int) string {
		return acc + fmt.Sprintf("$%d, ", i+2)
	})
	where += fmt.Sprintf(" AND id IN (%s)", strings.TrimSuffix(argsStr, ", "))
	args := utils.MapSlice(coupletIds, func(id string) any { return id })
	args = append([]any{songId}, args...)

	return r.queryAction(trx, where, args...)
}

func (r *Repository) FindOne(where string, args ...any) (*CoupletModel, res.Error) {

	couplet := CoupletModel{}

	err := r.db.Get(&couplet, where, args...)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, res.NewError(err, http.StatusInternalServerError)
	}

	return &couplet, nil
}

func (r *Repository) FindMany(where string, args ...any) ([]CoupletModel, res.Error) {
	couplets := []CoupletModel{}

	err := r.db.Select(&couplets, where, args...)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, res.NewError(err, http.StatusInternalServerError)
	}

	return couplets, nil
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

func (r *Repository) getExecutor(trx *sql.Tx) postgres.QueryFn {
	if trx != nil {
		return trx.Query
	}

	return r.db.Query

}

func (r *Repository) GetTrx() (*sql.Tx, error) {
	return r.db.Db.Begin()
}
