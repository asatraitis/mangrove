package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_init.go -package=mocks github.com/asatraitis/mangrove/internal/handler InitHandler
type InitHandler interface {
	home(http.ResponseWriter, *http.Request)
	initRegistration(http.ResponseWriter, *http.Request)
}
type initHandler struct {
	logger zerolog.Logger
	bll    bll.BLL

	initMux *http.ServeMux

	webauthn webauthn.WebAuthN
	config   config.Configs
}

type InitRegistrationRequest struct {
	RegistrationCode string `json:"registrationCode"`
}
type InitRegistrationResponse *protocol.CredentialCreation

func NewInitHandler(logger zerolog.Logger, bll bll.BLL, initMux *http.ServeMux, webauthn webauthn.WebAuthN, config config.Configs) InitHandler {
	logger = logger.With().Str("subcomponent", "InitHandler").Logger()
	h := &initHandler{
		logger:   logger,
		bll:      bll,
		initMux:  initMux,
		webauthn: webauthn,
		config:   config,
	}
	h.register()
	return h
}

func (ih *initHandler) register() {
	ih.initMux.Handle("GET /", http.FileServer(http.Dir("./dist/init")))
	ih.initMux.HandleFunc("POST /", ih.initRegistration)
}

func (ih *initHandler) home(w http.ResponseWriter, r *http.Request) {
	ih.logger.Info().Str("func", "home").Msg("GET /")
	w.WriteHeader(http.StatusOK)
}

func (ih *initHandler) initRegistration(w http.ResponseWriter, r *http.Request) {
	// TODO: Refacotr; too much logic in a handler - should be in BLL
	var reg InitRegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// regCode is hashed value
	regCode, err := ih.config.GetConfig(dal.CONFIG_INIT_SA_CODE)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ih.logger.Info().Msgf("OK: %s - %s", regCode, reg.RegistrationCode)

	w.WriteHeader(http.StatusOK)
	// response, err := ih.webauthn.BeginRegistration()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
