package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
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
	ih.initMux.HandleFunc("POST /v1/register", ih.initRegistration)
	ih.initMux.HandleFunc("POST /v1/register/finish", HandleWithMiddleware(ih.finishRegistration,
		[]MiddlewareFunc{
			ih.middleware.CsrfValidationMiddleware,
		},
	))
}

// TODO: remove
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

	// FIXME: this is useless as it stands. Need to sign an access token + device fingerprint or something else
	// Currently, it signs a random token that could be reused between sessions
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
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	err = ih.bll.User(ctx).RegisterSuperAdmin(&req)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	// set an auth token; should not error ever - same parsing done in BLL previously
	bUserID, err := base64.StdEncoding.DecodeString(req.UserID)
	if err != nil {
		ih.logger.Error().Msg("failed to parse user id")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(string(bUserID))
	if err != nil {
		ih.logger.Error().Msg("failed to parse user id")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	token, err := ih.bll.User(ctx).CreateToken(userUUID)
	if err != nil {
		ih.logger.Error().Msg("failed to create user token")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token.ID.String(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// Secure: true, // TODO: this needs to be set TRUE for prod
	})

	w.WriteHeader(http.StatusOK)
}

// TODO: move it to utils within handler package?
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
