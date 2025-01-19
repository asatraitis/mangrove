package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/handler/types"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
)

type Middleware interface {
	CsrfValidationMiddleware(HandlerFuncType) HandlerFuncType
	AuthValidationMiddleware(HandlerFuncType) HandlerFuncType
}
type middleware struct {
	vars *configs.EnvVariables
	bll  bll.BLL
}

func NewMiddleware(vars *configs.EnvVariables, bll bll.BLL) Middleware {
	return &middleware{
		vars: vars,
		bll:  bll,
	}
}
func (m *middleware) AuthValidationMiddleware(next HandlerFuncType) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie("auth_token")
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		if authToken.Value != "" {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		tokenID, err := uuid.Parse(authToken.Value)
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		user, err := m.bll.User(r.Context()).ValidateTokenAndGetUser(tokenID)
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		// TODO: keep in context?
		ctx := context.WithValue(r.Context(), types.REQ_CTX_KEY_USER_TOKEN, authToken.Value)
		ctx = context.WithValue(ctx, types.REQ_CTX_KEY_USER_ID, user.ID.String())
		ctx = context.WithValue(ctx, types.REQ_CTX_KEY_USER_ROLE, user.Role)
		ctx = context.WithValue(ctx, types.REQ_CTX_KEY_USER_STATUS, user.Status)

		next(w, r.WithContext(ctx))
	}
}
func (m *middleware) CsrfValidationMiddleware(next HandlerFuncType) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := r.Header.Get("X-CSRF-Token")
		csrfCookie, err := r.Cookie("csrf_token")
		if err != nil || csrfToken == "" {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		csrfParts := strings.Split(csrfCookie.Value, ".")
		if len(csrfParts) != 2 {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		if csrfToken != csrfParts[0] {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		hasher := utils.NewStandardCrypto([]byte(m.vars.MangroveSalt))
		err = hasher.VerifyToken(csrfToken, csrfParts[1])
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
