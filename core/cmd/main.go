package main

import (
	"effect-mobile/config"
	"effect-mobile/internal/couplet"
	"effect-mobile/internal/song"
	"effect-mobile/pkg/logger"
	"effect-mobile/pkg/middlewares"
	"effect-mobile/pkg/postgres"
	"effect-mobile/pkg/res"
	"effect-mobile/pkg/utils"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.GetLogger()
	db := postgres.NewPostgresDB(cfg.DSN, "postgres", logger)

	router := httprouter.New()

	applyHandlers(router, db, logger)

	server := configureServer(router, cfg)

	logger.Info(fmt.Sprintf("Server started on port: %v", cfg.ApiPort))

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func configureServer(router *httprouter.Router, cfg *config.Config) *http.Server {

	const TEN_SECOND = time.Second * 10

	return &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.ApiPort),
		ReadTimeout:  TEN_SECOND,
		WriteTimeout: TEN_SECOND,
		Handler:      router,
	}
}
func applyHandlers(router *httprouter.Router, db *postgres.PostgresDb, logger *slog.Logger) {

	middlewares := middlewares.Chain(middlewares.LogMid(logger), middlewares.CORS)

	//REPOSITORIES
	songRepo := song.NewRepository(logger, db)
	coupletRepo := couplet.NewRepository(db, logger)

	//SERVICES
	songService := song.NewService(songRepo, logger)
	coupletService := couplet.NewService(coupletRepo, songRepo, logger)

	//HANDLERS
	song.ApplyHandler(router, songService, middlewares, logger)
	couplet.ApplyHandler(router, coupletService, middlewares, logger)

	router.PanicHandler = panicHandler(logger)
	router.NotFound = http.HandlerFunc(notFoundHandler(logger))

}

func panicHandler(logger *slog.Logger) func(w http.ResponseWriter, r *http.Request, p interface{}) {
	return func(w http.ResponseWriter, r *http.Request, p interface{}) {
		msg := fmt.Sprintf("Something went wrong, err: %v", p)
		logger.Error("Panic: " + msg)
		res.NewErrorWithMessage(msg, 500).ResHttp(w)
	}
}

func notFoundHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.LogUnhandlingRequest(r, logger)
		msg := fmt.Sprintf("PATH: %v by METHOD: %v not available", r.URL.Path, r.Method)
		err := res.NewErrorWithMessage(msg, http.StatusNotFound)
		res.HttpJsonRes(w, err, err.Status)
	}
}
