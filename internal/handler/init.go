package handler

import (
	"net/http"

	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_init.go -package=mocks github.com/asatraitis/mangrove/internal/handler InitHandler
type InitHandler interface {
	home(http.ResponseWriter, *http.Request)
}
type initHandler struct {
	logger zerolog.Logger
	bll    bll.BLL

	initMux *http.ServeMux
}

func NewInitHandler(logger zerolog.Logger, bll bll.BLL, initMux *http.ServeMux) InitHandler {
	logger = logger.With().Str("subcomponent", "InitHandler").Logger()
	h := &initHandler{
		logger:  logger,
		bll:     bll,
		initMux: initMux,
	}
	h.register()
	return h
}

func (ih *initHandler) register() {
	ih.initMux.Handle("GET /", http.FileServer(http.Dir("./dist/init")))
}

func (ih *initHandler) home(w http.ResponseWriter, r *http.Request) {
	ih.logger.Info().Str("func", "home").Msg("GET /")
	w.WriteHeader(http.StatusOK)
}
