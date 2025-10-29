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
	Cache       *cache.Storage
	Store       *store.Storage
	Translator  translate.Translator
	Transporter *transport.Transporter
	Logger      *zap.SugaredLogger
}

type HandlerConfig struct {
	Config      config.APIConfig
	Cache       *cache.Storage
	Store       *store.Storage
	Translator  translate.Translator
	Transporter *transport.Transporter
	Logger      *zap.SugaredLogger
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		Config:      cfg.Config,
		Cache:       cfg.Cache,
		Store:       cfg.Store,
		Translator:  cfg.Translator,
		Transporter: cfg.Transporter,
		Logger:      cfg.Logger,
	}
}
