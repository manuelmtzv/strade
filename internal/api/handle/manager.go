package handle

import (
	"strade/internal/api/transport"
	"strade/internal/cache"
	"strade/internal/config"
	"strade/internal/store"
	"strade/internal/translate"

	"go.uber.org/zap"
)

type Handler struct {
	Config      config.APIConfig
	Cache       cache.Storage
	Store       store.Storage
	Translator  translate.Translator
	Transporter *transport.Transporter
	Logger      *zap.SugaredLogger
}

func NewHandler(config config.APIConfig, cache cache.Storage, store store.Storage, translator translate.Translator, transporter *transport.Transporter, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		Config:      config,
		Cache:       cache,
		Store:       store,
		Translator:  translator,
		Transporter: transporter,
		Logger:      logger,
	}
}
