package song

import (
	"effect-mobile/pkg/middlewares"
	"effect-mobile/pkg/req"
	"effect-mobile/pkg/res"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	createSong   = "/songs"
	listingSongs = createSong
	getSong      = "/songs/:id"
	deleteSong   = getSong
	updateSong   = getSong
	info         = "/info"
)

type handler struct {
	logger      *slog.Logger
	songService *Service
}

func ApplyHandler(router *httprouter.Router, songService *Service, middleware middlewares.Middleware, logger *slog.Logger) {

	handler := handler{
		songService: songService,
		logger:      logger,
	}

	router.POST(createSong, middleware(handler.Create))
	router.PATCH(updateSong, middleware(handler.Update))
	router.DELETE(deleteSong, middleware(handler.Delete))
	router.GET(getSong, middleware(handler.GetOne))
	router.GET(listingSongs, middleware(handler.Listing))
	router.GET(info, middleware(handler.Info))

	logHandlerRouters(logger)

}

func (h *handler) Info(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	query := r.URL.Query()

	song := query.Get("song")
	group := query.Get("group")

	if song == "" {
		msg := "Song in query must be passed!"
		h.logger.Error(msg)
		res.NewErrorWithMessage(msg, 400).ResHttp(w)
		return
	}

	if group == "" {
		msg := "Group in query must be passed!"
		h.logger.Error(msg)
		res.NewErrorWithMessage(msg, 400).ResHttp(w)
		return
	}

	result, err := h.songService.Info(group, song)

	if err != nil {
		err.ResHttp(w)
		return
	}

	res.HttpJsonRes(w, result.MapToShow(), 200)

}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")

	if err := h.songService.Delete(id); err != nil {
		err.ResHttp(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info("Song was deleted!")

}

func (h *handler) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()

	id := p.ByName("id")

	body, err := req.HandleBody[SongUpdateDto](r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		err.ResHttp(w)
		return
	}

	err = h.songService.Update(id, &body)

	if err != nil {
		err.ResHttp(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()

	body, err := req.HandleBody[SongCreateDto](r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		err.ResHttp(w)
		return
	}

	resData, err := h.songService.Create(&body)

	if err != nil {
		err.ResHttp(w)
		return
	}

	h.logger.Info("New song was created!")
	res.HttpJsonRes(w, resData, 201)

}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()

	id := p.ByName("id")

	if uuid.Validate(id) != nil {
		res.
			NewErrorWithMessage("Invalid songId!", http.StatusBadRequest).
			ResHttp(w)
		return
	}

	song, err := h.songService.FindOne(id)

	if err != nil {
		err.ResHttp(w)
		return
	}

	res.HttpJsonRes(w, song.MapToShow(), http.StatusOK)

}
func (h *handler) Listing(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()

	pagParams := res.MapQueryToPagParams(r.URL.Query())

	result, err := h.songService.Listing(pagParams)

	if err != nil {
		err.ResHttp(w)
		return
	}

	res.HttpJsonRes(w, result, 200)

}

func logHandlerRouters(logger *slog.Logger) {
	logger.Info("Handler: GET /info")
	logger.Info("Song handler: POST " + createSong)
	logger.Info("Song handler: PATCH " + updateSong)
	logger.Info("Song handler: DELETE " + deleteSong)
	logger.Info("Song handler: GET " + getSong)
	logger.Info("Song handler: GET " + listingSongs)
	logger.Info("Song handler was applied!")

}
