package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/utils"
)

//go:generate mockgen -destination=./mocks/mock_init.go -package=mocks github.com/asatraitis/mangrove/internal/handler InitHandler
type InitHandler interface {
	home(http.ResponseWriter, *http.Request)
	initRegistration(http.ResponseWriter, *http.Request)
	finishRegistration(http.ResponseWriter, *http.Request)
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
	ih.initMux.HandleFunc("POST /finish", ih.csrfValidationMiddleware(ih.finishRegistration))
}

func (ih *initHandler) home(w http.ResponseWriter, r *http.Request) {
	ih.logger.Info().Str("func", "home").Msg("GET /")
	w.WriteHeader(http.StatusOK)
}

func (ih *initHandler) initRegistration(w http.ResponseWriter, r *http.Request) {
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

	userRegCreds, err := ih.bll.User(ctx).CreateUserSession()
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "failed to create registration credentials",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}
	res.PublicKey = userRegCreds.Response

	hasher := utils.NewStandardCrypto([]byte(ih.vars.MangroveSalt))
	token, sig, err := hasher.GenerateTokenHMAC()
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "failed to create a signiture",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token + "." + sig,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// Secure: true, // TODO: this needs to be set TRUE for prod
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.InitRegistrationResponse]{Response: &res})
}

func (ih *initHandler) finishRegistration(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context = context.Background()

	var req dto.FinishRegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	err = ih.bll.User(ctx).RegisterSuperAdmin(&req)
	if err != nil {
		sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
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

// TODO: this middleware needs to be made reusable in other handler structs; depends on a salt
// that is part of the env vars dependency
func (ih *initHandler) csrfValidationMiddleware(next HandlerFuncType) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := r.Header.Get("X-CSRF-Token")
		csrfCookie, err := r.Cookie("csrf_token")
		if err != nil || csrfToken == "" {
			sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		csrfParts := strings.Split(csrfCookie.Value, ".")
		if len(csrfParts) != 2 {
			sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		if csrfToken != csrfParts[0] {
			sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		hasher := utils.NewStandardCrypto([]byte(ih.vars.MangroveSalt))
		err = hasher.VerifyToken(csrfToken, csrfParts[1])
		if err != nil {
			sendErrResponse[dto.InitRegistrationResponse](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
