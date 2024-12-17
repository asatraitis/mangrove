package handler

import (
	"net/http"

	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_handler.go -package=mocks github.com/asatraitis/mangrove/internal/handler Handler
type Handler interface {
	Init(*http.ServeMux) InitHandler
}
type handler struct {
	logger   zerolog.Logger
	bll      bll.BLL
	config   config.Configs
	webauthn webauthn.WebAuthN
}

func NewHandler(logger zerolog.Logger, bll bll.BLL, webauthn webauthn.WebAuthN, config config.Configs) Handler {
	logger = logger.With().Str("component", "Handler").Logger()
	return &handler{
		logger:   logger,
		bll:      bll,
		webauthn: webauthn,
		config:   config,
	}
}

func (h *handler) Init(mux *http.ServeMux) InitHandler {
	return NewInitHandler(h.logger, h.bll, mux, h.webauthn, h.config)
}
