package postgres

import (
	"database/sql"
	"effect-mobile/pkg/utils"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type QueryFn = func(query string, args ...any) (*sql.Rows, error)

type PostgresDb struct {
	Db     *sqlx.DB
	logger *slog.Logger
}

func (db *PostgresDb) Get(dest interface{}, query string, args ...interface{}) error {
	db.logger.Info(query, args...)
	return db.Db.Get(dest, query, args...)
}

func (db *PostgresDb) Query(query string, args ...any) (*sql.Rows, error) {
	db.logger.Info(query, args...)
	return db.Db.Query(query, args...)
}

func (db *PostgresDb) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	db.logger.Info(query, args...)
	return db.Db.QueryRowx(query, args...)
}

func (db *PostgresDb) Select(dest interface{}, query string, args ...interface{}) error {
	db.logger.Info(query, args...)
	return db.Db.Select(dest, query, args...)
}

func NewPostgresDB(dsn, dialect string, logger *slog.Logger) *PostgresDb {

	db, err := sqlx.Open(dialect, dsn)
	if err != nil {
		logger.Error("Error in connecting to DB")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Ping to DB responsed error")
		logger.Error(err.Error())
		panic(err)
	}
	logger.Info("Successfully connected to DB")

	poolConnCount := utils.ParseStringToIntOrDefault(os.Getenv("CONNECTIONS_COUNT"), 100)
	db.SetMaxOpenConns(poolConnCount)
	logger.Info(fmt.Sprintf("Pool connections  count is %v", poolConnCount))

	return &PostgresDb{
		Db:     db,
		logger: logger,
	}
}
