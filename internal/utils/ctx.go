package utils

import (
	"context"
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type LocalizerKey string

const (
	LocalizerCtx LocalizerKey = "localizer"
)

func GetLocalizerFromContext(ctx context.Context) (*i18n.Localizer, error) {
	localizer, ok := ctx.Value(LocalizerCtx).(*i18n.Localizer)
	if !ok {
		return nil, errors.New("localizer not found in context")
	}
	return localizer, nil
}
