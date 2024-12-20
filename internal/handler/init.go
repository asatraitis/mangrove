package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/dto"
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
	// TODO: Need to add a middleware to add a signed token to a cookie
	// csrf_token = "<raw_token>|<signature>"
	// on requests validate that <raw_token> in X-CSRF-Token header using the <signature>
	var req dto.InitRegistrationRequest
	var ctx context.Context = context.Background()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = ih.bll.Config(ctx).ValidateRegistrationCode(req.RegistrationCode)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userRegCreds, csrfToken, err := ih.bll.User(ctx).CreateUserSession()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: csrfToken,
		Path:  "/",
		// Secure: true, // TODO: this needs to be set TRUE for prod
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := dto.Response[dto.InitRegistrationResponse]{Response: &dto.InitRegistrationResponse{
		PublicKey: userRegCreds.Response,
	}}
	json.NewEncoder(w).Encode(res)

	// regCode is hashed value
	// regCode, err := ih.appConfig.GetConfig(dal.CONFIG_INIT_SA_CODE)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// ih.logger.Info().Msgf("OK: %s - %s", regCode, reg.RegistrationCode)

	// w.WriteHeader(http.StatusOK)
	// response, err := ih.webauthn.BeginRegistration()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
