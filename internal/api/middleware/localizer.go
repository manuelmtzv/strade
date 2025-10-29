package middleware

import (
	"net/http"
	"strade/internal/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Localizer(bundle *i18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(bundle, accept)
			ctx := utils.SetLocalizerInContext(r.Context(), localizer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
