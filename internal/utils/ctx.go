package utils

import (
	"context"
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type contextKey struct{}

var localizerKey = contextKey{}

func SetLocalizerInContext(ctx context.Context, localizer *i18n.Localizer) context.Context {
	return context.WithValue(ctx, localizerKey, localizer)
}

func GetLocalizerFromContext(ctx context.Context) (*i18n.Localizer, error) {
	localizer, ok := ctx.Value(localizerKey).(*i18n.Localizer)
	if !ok {
		return nil, errors.New("localizer not found in context")
	}
	return localizer, nil
}
