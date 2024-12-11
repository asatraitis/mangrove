package service

import (
	"net/http"

	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/handler"
	"github.com/rs/zerolog"
)

type Router interface {
	http.Handler
}
type router struct {
	logger  zerolog.Logger
	configs Configs
	bll     bll.BLL

	initMux *http.ServeMux
	mainMux *http.ServeMux
}

func NewRouter(logger zerolog.Logger, configs Configs, bll bll.BLL) Router {
	logger = logger.With().Str("component", "Router").Logger()
	ro := &router{
		logger:  logger,
		configs: configs,
		bll:     bll,

		initMux: http.NewServeMux(),
		mainMux: http.NewServeMux(),
	}
	ro.register()
	return ro
}

func (ro *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	instanceReady, err := ro.configs.GetConfig(dal.CONFIG_INSTANCE_READY)
	if err != nil {
		ro.logger.Err(err).Str("func", "ServeHTTP").Msg("failed to route request: config error")
		w.WriteHeader(http.StatusInternalServerError)
	}
	if instanceReady == "true" {
		ro.mainMux.ServeHTTP(w, r)
	} else {
		ro.initMux.ServeHTTP(w, r)
	}
}

func (ro *router) register() {
	handler.NewInitHandler(ro.logger, ro.bll, ro.initMux)
}
