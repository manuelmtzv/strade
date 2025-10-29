package translate

import (
	"context"
	"errors"
	"strade/internal/utils"
	"text/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type Translator interface {
	GetLocaleMessage(context.Context, string, map[string]any) (string, error)
	GetMessageOrDefault(context.Context, string, string, map[string]any) string
	GetLocaleError(context.Context, string, map[string]any) error
	GetLocaleHtmlTemplate(context.Context, string) (*template.Template, error)
}

type DefaultTranslator struct {
	logger *zap.SugaredLogger
}

func NewDefaultTranslator(logger *zap.SugaredLogger) *DefaultTranslator {
	return &DefaultTranslator{
		logger: logger,
	}
}

func (t *DefaultTranslator) GetLocaleMessage(ctx context.Context, messageID string, templateData map[string]any) (string, error) {
	l, err := utils.GetLocalizerFromContext(ctx)
	if err != nil {
		t.logger.Error("failed to get localizer from context",
			zap.String("messageID", messageID),
			zap.Any("templateData", templateData),
			zap.Error(err),
		)
		return "", err
	}

	return l.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

func (t *DefaultTranslator) GetMessageOrDefault(ctx context.Context, messageID, defaultMessage string, templateData map[string]any) string {
	msg, err := t.GetLocaleMessage(ctx, messageID, templateData)

	if err != nil {
		return defaultMessage
	}
	return msg
}

func (t *DefaultTranslator) GetLocaleError(ctx context.Context, messageID string, templateData map[string]any) error {
	msg, err := t.GetLocaleMessage(ctx, messageID, templateData)
	if err != nil {
		return err
	}
	return errors.New(msg)
}

func (t *DefaultTranslator) GetLocaleHtmlTemplate(ctx context.Context, pathID string) (*template.Template, error) {
	path, err := t.GetLocaleMessage(ctx, pathID, nil)
	if err != nil {
		return nil, err
	}

	return template.ParseFiles(path)
}
