package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
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
	h.mux.HandleFunc("GET /v1/me", h.middleware.AuthValidationMiddleware(h.middleware.CsrfValidationMiddleware(h.me)))
	h.mux.HandleFunc("POST /v1/login", h.initLogin)
	h.mux.Handle("GET /", http.FileServer(http.Dir("./dist/main")))
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
	userID, err := GetUserIdFromCtx(ctx)
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

	hasher := utils.NewStandardCrypto([]byte(h.vars.MangroveSalt))
	token, sig, err := hasher.GenerateTokenHMAC()
	if err != nil {
		sendErrResponse[any](w, &dto.ResponseError{
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

	json.NewEncoder(w).Encode(dto.Response[dto.InitLoginResponse]{Response: &res})
}
