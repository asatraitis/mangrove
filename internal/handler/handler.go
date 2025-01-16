package handler

import (
	"net/http"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_handler.go -package=mocks github.com/asatraitis/mangrove/internal/handler Handler
type Handler interface {
	Init(*http.ServeMux) InitHandler
	Main(*http.ServeMux) MainHandler
}
type BaseHandler struct {
	logger     zerolog.Logger
	vars       *configs.EnvVariables
	appConfig  config.Configs // TODO: Do we need it in the handler?
	bll        bll.BLL
	middleware Middleware
}
type handler struct {
	*BaseHandler
}
type HandlerFuncType func(w http.ResponseWriter, r *http.Request)

func NewHandler(logger zerolog.Logger, bll bll.BLL, vars *configs.EnvVariables, appConfig config.Configs) Handler {
	logger = logger.With().Str("component", "Handler").Logger()
	return &handler{
		BaseHandler: &BaseHandler{
			logger:     logger,
			vars:       vars,
			appConfig:  appConfig,
			bll:        bll,
			middleware: NewMiddleware(vars),
		},
	}
}

func (h *handler) Init(mux *http.ServeMux) InitHandler {
	return NewInitHandler(h.BaseHandler, mux)
}
func (h *handler) Main(mux *http.ServeMux) MainHandler {
	return NewMainHandler(h.BaseHandler, mux)
}
