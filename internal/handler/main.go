package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
)

type MainHandler interface{}
type mainHandler struct {
	*BaseHandler

	mux *http.ServeMux
}

func NewMainHandler(baseHandler *BaseHandler, mux *http.ServeMux) MainHandler {
	h := &mainHandler{
		BaseHandler: baseHandler,
		mux:         mux,
	}
	h.logger = h.logger.With().Str("subcomponent", "MainHandler").Logger()
	h.register()
	return h
}

func (h *mainHandler) register() {
	h.mux.HandleFunc("GET /v1/me", HandleWithMiddleware(h.me,
		[]MiddlewareFunc{
			h.middleware.CsrfValidationMiddleware,
			h.middleware.AuthValidationMiddleware,
		},
	))
	h.mux.HandleFunc("POST /v1/login", h.initLogin)
	h.mux.HandleFunc("POST /v1/login/finish", HandleWithMiddleware(h.finishLogin,
		[]MiddlewareFunc{
			h.middleware.CsrfValidationMiddleware,
		},
	))
	h.mux.Handle("GET /", http.FileServer(http.Dir("./dist/main")))
	h.mux.HandleFunc("GET /v1/clients", HandleWithMiddleware(
		h.clients,
		[]MiddlewareFunc{
			h.middleware.CsrfValidationMiddleware,
			h.middleware.AuthValidationMiddleware,
			h.middleware.UserStatusValidation,
			h.middleware.UserRoleSuperadmin,
		},
	))

	h.mux.HandleFunc("POST /v1/clients", HandleWithMiddleware(
		h.createClient,
		[]MiddlewareFunc{
			h.middleware.CsrfValidationMiddleware,
			h.middleware.AuthValidationMiddleware,
			h.middleware.UserStatusValidation,
			h.middleware.UserRoleSuperadmin,
		},
	))

}
func (h *mainHandler) clientRouting() http.Handler {
	const fsPath = "./dist/main"
	fs := http.FileServer(http.Dir(fsPath))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the requested file exists then return if; otherwise return index.html (fileserver default page)
		if r.URL.Path != "/" {
			fullPath := fsPath + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, r)
	})
}
func (h *mainHandler) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := &dto.MeResponse{}

	// get user UUID
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		h.logger.Err(err).Msg("failed to get userID from context")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to retrieve user id",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.bll.User(ctx).GetUserByID(userID)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to retrieve user",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	// models.User -> dto.MeResponse
	res, err = typeconv.ConvertUserToMeResponse(user)
	if err != nil {
		h.logger.Err(err)
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to typeconv",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.MeResponse]{Response: res})
}

func (h *mainHandler) initLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.InitLoginRequest
	var res dto.InitLoginResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	creds, sessionKey, err := h.bll.User(ctx).InitLogin(req.Username)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to create login credentials",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}
	res.PublicKey = creds
	res.SessionKey = sessionKey

	h.setCsrfCookies(w, r, "")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.InitLoginResponse]{Response: &res})
}

func (h *mainHandler) finishLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.FinishLoginRequest
	var res *dto.MeResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	res, err = h.bll.User(ctx).FinishLogin(&req)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to login",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(res.ID)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to parse uuid",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	token, err := h.bll.User(ctx).CreateToken(userID)
	if err != nil {
		h.logger.Error().Msg("failed to create user token")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: err.Error(),
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	h.setCsrfCookies(w, r, token.ID.String())

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token.ID.String(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// Secure: true, // TODO: this needs to be set TRUE for prod
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.MeResponse]{Response: res})
}
func (h *mainHandler) clients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userClients, err := h.bll.Client(ctx).GetUserClients()
	if err != nil {
		h.logger.Err(err).Msg("failed to get userID from context")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to get user clients",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.Response[dto.UserClientsResponse]{Response: &userClients})

}
func (h *mainHandler) createClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateClientRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Err(err).Msg("failed to decode payload")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "invalid request body",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	res, err := h.bll.Client(ctx).Create(req)
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to create client",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(dto.Response[dto.CreateClientResponse]{Response: res})
}

// TODO: Consolidate with init handler (dupe)
func (h *mainHandler) setCsrfCookies(w http.ResponseWriter, r *http.Request, authToken string) {
	h.logger.Info().Msg("setting CSRF Cookies")
	if r == nil || w == nil {
		h.logger.Error().Msg("missing request or response; failed to set csrf cookies")
		return
	}

	// validating IP might be harder to use with a proxy - omitting for now
	// randomUUID + authToken + userAgent + acceptHeader + acceptLanguageHeader
	csrfToken := uuid.NewString()
	ip := getReqIP(r)
	data := csrfToken + authToken + r.UserAgent() + r.Header.Get("Accept") + r.Header.Get("Accept-Language")
	h.logger.Info().Str("userAgent", r.UserAgent()).Str("RemoteAdds", ip).Msg("data to sign")
	hasher := utils.NewStandardCrypto([]byte(h.vars.MangroveSalt))
	signature := hasher.GenerateTokenHMAC(data)

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   h.vars.MangroveEnv == configs.PROD,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_sig",
		Value:    signature,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   h.vars.MangroveEnv == configs.PROD,
		HttpOnly: true,
	})
}
