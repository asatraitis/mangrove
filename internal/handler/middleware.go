package handler

import (
	"context"
	"net/http"
	"slices"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/handler/types"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MiddlewareFunc func(HandlerFuncType) HandlerFuncType
type Middleware interface {
	CsrfValidationMiddleware(HandlerFuncType) HandlerFuncType
	AuthValidationMiddleware(HandlerFuncType) HandlerFuncType
	UserStatusValidation(HandlerFuncType) HandlerFuncType
	UserRoleSuperadmin(HandlerFuncType) HandlerFuncType
	UserRoleUser(next HandlerFuncType) HandlerFuncType
}
type middleware struct {
	vars   *configs.EnvVariables
	bll    bll.BLL
	logger zerolog.Logger
}

func NewMiddleware(vars *configs.EnvVariables, bll bll.BLL, logger zerolog.Logger) Middleware {
	return &middleware{
		vars:   vars,
		bll:    bll,
		logger: logger,
	}
}
func HandleWithMiddleware(handler HandlerFuncType, middleware []MiddlewareFunc) HandlerFuncType {
	var final HandlerFuncType = handler
	for _, mw := range slices.Backward(middleware) {
		final = mw(final)
	}

	return final
}
func (m *middleware) AuthValidationMiddleware(next HandlerFuncType) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info().Msg("Validating authentication status")
		authToken, err := r.Cookie("auth_token")
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token: no auth token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		if authToken.Value == "" {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token: no auth token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		tokenID, err := uuid.Parse(authToken.Value)
		if err != nil {
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token: bad token",
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
		m.logger.Info().Msg("Validating csrf token")
		csrfToken := r.Header.Get("X-CSRF-Token")
		csrfCookie, err := r.Cookie("csrf_token")
		if err != nil || csrfToken == "" {
			m.logger.Err(err).Msg("failed to get csrf_token cookie")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		if csrfToken != csrfCookie.Value {
			m.logger.Err(err).Msg("csrf token did not match cookie")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		authCookie, _ := r.Cookie("auth_token")
		var authToken string
		if authCookie != nil {
			authToken = authCookie.Value
		}

		csrfSigCookie, err := r.Cookie("csrf_sig")
		if err != nil {
			m.logger.Err(err).Msg("failed to get csrf_sig cookie")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		hasher := utils.NewStandardCrypto([]byte(m.vars.MangroveSalt))
		// csrfToken + authToken + userAgent + acceptHeader + acceptLanguageHeader
		err = hasher.VerifyToken(csrfToken+authToken+r.UserAgent()+r.Header.Get("Accept")+r.Header.Get("Accept-Language"), csrfSigCookie.Value)
		if err != nil {
			m.logger.Err(err).Str("data", csrfToken+authToken+r.UserAgent()+r.Header.Get("Accept")+r.Header.Get("Accept-Language")).Str("sig", csrfSigCookie.Value).Msg("failed to verify csrf")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate token",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
func (m *middleware) UserStatusValidation(next HandlerFuncType) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info().Msg("Validating user status")
		status, err := utils.GetUserStatusFromCtx(r.Context())
		if err != nil {
			m.logger.Err(err).Str("status", string(status)).Msg("failed to get user status from context")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate user status",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		if status != models.USER_STATUS_ACTIVE {
			m.logger.Err(err).Str("status", string(status)).Msg("user status is not active")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "user is not active",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
func (m *middleware) UserRoleSuperadmin(next HandlerFuncType) HandlerFuncType {
	return m.userRoleValidation(next, models.USER_ROLE_SUPERUSER)
}
func (m *middleware) UserRoleUser(next HandlerFuncType) HandlerFuncType {
	return m.userRoleValidation(next, models.USER_ROLE_USER)
}

func (m *middleware) userRoleValidation(next HandlerFuncType, role models.UserRole) HandlerFuncType {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info().Msg("Validating user status")
		userRole, err := utils.GetUserRoleFromCtx(r.Context())
		if err != nil {
			m.logger.Err(err).Str("role", string(userRole)).Msg("failed to get user status from context")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "failed to validate user status",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusBadRequest)
			return
		}
		if userRole != role {
			m.logger.Err(err).Str("role", string(userRole)).Str("allowedRole", string(role)).Msg("user role was not the role validate against")
			sendErrResponse[any](w, &dto.ResponseError{
				Message: "not allowed",
				Code:    "ERROR_CODE_TBD",
			}, http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
