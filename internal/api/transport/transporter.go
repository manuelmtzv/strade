package transport

import (
	"strade/internal/translate"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Transporter struct {
	Validate   *validator.Validate
	logger     *zap.SugaredLogger
	translator translate.Translator
}

func NewTransporter(v *validator.Validate, l *zap.SugaredLogger, t translate.Translator) *Transporter {
	return &Transporter{
		Validate:   v,
		logger:     l,
		translator: t,
	}
}
