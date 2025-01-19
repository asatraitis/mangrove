package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
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
	h.mux.Handle("GET /", http.FileServer(http.Dir("./dist/main")))
	h.mux.HandleFunc("GET /v1/me", h.middleware.AuthValidationMiddleware(h.me))
}

func (h *mainHandler) me(w http.ResponseWriter, r *http.Request) {
	res := &dto.MeResponse{}

	// get user UUID
	userID, err := GetUserIdFromCtx(r.Context())
	if err != nil {
		h.logger.Err(err).Msg("failed to get userID from context")
		sendErrResponse[any](w, &dto.ResponseError{
			Message: "failed to retrieve user id",
			Code:    "ERROR_CODE_TBD",
		}, http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.bll.User(r.Context()).GetUserByID(userID)
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
