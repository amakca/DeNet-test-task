package middlewares

import (
	"context"
	"denet-test-task/internal/api/v1/apierrs"
	"denet-test-task/internal/services/auth"
	"denet-test-task/pkg/logctx"
	"net/http"
	"strings"
)

type ctxKey string

const (
	userIdCtx ctxKey = "userId"
)

type AuthMiddleware struct {
	AuthService auth.Auth
}

func (h *AuthMiddleware) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logctx.FromContext(r.Context())

		token, ok := bearerToken(r)
		if !ok {
			log.Warn("AuthMiddleware.UserIdentity: bearerToken", "error", apierrs.ErrInvalidAuthHeader)
			apierrs.NewErrorResponseHTTP(w, http.StatusUnauthorized, apierrs.ErrInvalidAuthHeader.Error())
			return
		}

		userId, err := h.AuthService.ParseToken(token)
		if err != nil {
			log.Warn("AuthMiddleware.UserIdentity: ParseToken", "err", err)
			apierrs.NewErrorResponseHTTP(w, http.StatusUnauthorized, apierrs.ErrCannotParseToken.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
