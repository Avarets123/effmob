package couplet

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
	addCouplets     = "/songs/:id/couplets"
	listingCouplets = addCouplets
	deleteCouplets  = addCouplets
	rewriteCouplets = addCouplets
)

type handler struct {
	coupletService *Service
	logger         *slog.Logger
}

func ApplyHandler(router *httprouter.Router, coupletService *Service, middleware middlewares.Middleware, logger *slog.Logger) {
	handler := handler{
		coupletService: coupletService,
		logger:         logger,
	}

	router.GET(listingCouplets, middleware(handler.Listing))
	router.PATCH(addCouplets, middleware(handler.AddCouplets))
	router.PUT(rewriteCouplets, middleware(handler.RewriteCouplets))
	router.DELETE(deleteCouplets, middleware(handler.Delete))

	logHandlerRouter(logger)
}

func (h *handler) AddCouplets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	songId := p.ByName("id")
	if uuid.Validate(songId) != nil {
		res.NewErrorWithMessage("Invalid uuid", 400).ResHttp(w)
		return
	}

	body, err := req.HandleBody[CoupletCreateDto](r.Body)
	if err != nil {
		err.ResHttp(w)
		return
	}

	if len(body.Couplets) == 0 {
		msg := "Couplets text must be passed!"
		h.logger.Error(msg)
		res.NewErrorWithMessage(msg, 422).ResHttp(w)
		return
	}

	err = h.coupletService.Create(nil, songId, body, true)
	if err != nil {
		err.ResHttp(w)
		return
	}

	w.WriteHeader(201)
	h.logger.Info("New couplets was created!")

}
func (h *handler) RewriteCouplets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	body, err := req.HandleBody[CoupletCreateDto](r.Body)
	songId := p.ByName("id")

	if uuid.Validate(songId) != nil {
		res.NewErrorWithMessage("Invalid uuid", 400).ResHttp(w)
		return
	}

	if err != nil {
		err.ResHttp(w)
		return
	}

	err = h.coupletService.RewriteCouplets(songId, body)

	if err != nil {
		err.ResHttp(w)
		return
	}

	w.WriteHeader(201)
	h.logger.Info("New couplets was rewrited!")
}
func (h *handler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	body, err := req.HandleBody[CoupletDeleteDto](r.Body)
	songId := p.ByName("id")

	if uuid.Validate(songId) != nil {
		res.NewErrorWithMessage("Invalid uuid", 400).ResHttp(w)
		return
	}

	if err != nil {
		err.ResHttp(w)
		return
	}

	err = h.coupletService.Delete(nil, songId, body.CoupletsIds)
	if err != nil {
		err.ResHttp(w)
		return
	}

	w.WriteHeader(204)
	h.logger.Info("Couplets was deleted!")

}
func (h *handler) Listing(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	songId := p.ByName("id")
	if uuid.Validate(songId) != nil {
		res.NewErrorWithMessage("Invalid uuid", 400).ResHttp(w)
		return
	}

	pagParams := res.MapQueryToPagParams(r.URL.Query())
	pagParams.SortField = "couplet_num"
	pagParams.SortDir = "ASC"

	result, err := h.coupletService.Listing(songId, pagParams)
	if err != nil {
		err.ResHttp(w)
		return
	}

	res.HttpJsonRes(w, result, 200)

}

func logHandlerRouter(logger *slog.Logger) {
	logger.Info("Couplets handler: PATCH " + addCouplets)
	logger.Info("Couplets handler: PUT " + rewriteCouplets)
	logger.Info("Couplets handler: GET " + listingCouplets)
	logger.Info("Couplets handler: DELETE " + deleteCouplets)
	logger.Info("Couplets handler was applied!")
}
