package router

import (
	"net/http"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/handler"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_router.go -package=mocks github.com/asatraitis/mangrove/internal/service/router Router
type Router interface {
	http.Handler

	register()
}
type router struct {
	logger  zerolog.Logger
	configs config.Configs
	handler handler.Handler

	initMux *http.ServeMux
	mainMux *http.ServeMux
}

func NewRouter(logger zerolog.Logger, configs config.Configs, handler handler.Handler) Router {
	logger = logger.With().Str("component", "Router").Logger()
	ro := &router{
		logger:  logger,
		configs: configs,
		handler: handler,

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
	ro.handler.Init(ro.initMux)
}
