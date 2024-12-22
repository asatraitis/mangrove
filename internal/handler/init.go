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
	var ctx context.Context = context.Background()

	var req dto.InitRegistrationRequest
	var res dto.InitRegistrationResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	err = ih.bll.Config(ctx).ValidateRegistrationCode(req.RegistrationCode)
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "invalid registration code",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	userRegCreds, csrfToken, err := ih.bll.User(ctx).CreateUserSession()
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "failed to create registration credentials",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}
	res.PublicKey = userRegCreds.Response

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // TODO: this needs to be set TRUE for prod
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.InitRegistrationResponse]{Response: &res})
}

func sendErrResponse[T any](w http.ResponseWriter, err *dto.ResponseError, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(
		dto.Response[T]{
			Response: nil,
			Error:    err,
		},
	)
}
