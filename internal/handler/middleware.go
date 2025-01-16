package handler

import (
	"net/http"
	"strings"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/utils"
)

type Middleware interface {
	CsrfValidationMiddleware(HandlerFuncType) HandlerFuncType
}
type middleware struct {
	vars *configs.EnvVariables
}

func NewMiddleware(vars *configs.EnvVariables) Middleware {
	return &middleware{
		vars: vars,
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
