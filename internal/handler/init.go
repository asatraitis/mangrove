package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/go-webauthn/webauthn/protocol"
)

//go:generate mockgen -destination=./mocks/mock_init.go -package=mocks github.com/asatraitis/mangrove/internal/handler InitHandler
type InitHandler interface {
	home(http.ResponseWriter, *http.Request)
	initRegistration(http.ResponseWriter, *http.Request)
}
type initHandler struct {
	*BaseHandler

	initMux *http.ServeMux
}

type InitRegistrationRequest struct {
	RegistrationCode string `json:"registrationCode"`
}
type InitRegistrationResponse *protocol.CredentialCreation

func NewInitHandler(baseHandler *BaseHandler, initMux *http.ServeMux) InitHandler {
	h := &initHandler{
		BaseHandler: baseHandler,
		initMux:     initMux,
	}
	h.logger = h.logger.With().Str("subcomponent", "InitHandler").Logger()
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
	regCode, err := ih.appConfig.GetConfig(dal.CONFIG_INIT_SA_CODE)
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
