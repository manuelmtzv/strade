package main

import (
	"database/sql"
	"strade/internal/api"
	"strade/internal/api/handle"
	"strade/internal/api/transport"
	"strade/internal/cache"
	"strade/internal/db"
	"strade/internal/env"
	"strade/internal/store"
	"strade/internal/translate"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	if err := env.Load(); err != nil {
		logger.Panicf("Failed to load env: %v", err)
	}

	cfg := getConfig()

	dbConn, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Panicf("Database connection failed: %v", err)
	}
	defer func(dbConn *sql.DB) {
		_ = dbConn.Close()
	}(dbConn)
	logger.Info("Database connection established")

	var redisClient *redis.Client
	if cfg.Cache.Redis.Enabled {
		redisClient = cache.NewRedisClient(cfg.Cache.Redis.Addr, cfg.Cache.Redis.PW, cfg.Cache.Redis.DB)
		defer func() {
			_ = redisClient.Close()
		}()

		if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
			logger.Panicf("Redis connection failed: %v", err)
		}

		logger.Info("Redis connection established")
	}

	storage := store.NewStorage(dbConn)
	cacheStorage := cache.NewStorage(redisClient)
	defaultTranslator := translate.NewDefaultTranslator(logger)
	validate := validator.New(validator.WithRequiredStructEnabled())
	transporter := transport.NewTransporter(validate, logger, defaultTranslator)

	handler := handle.NewHandler(cfg, cacheStorage, storage, defaultTranslator, transporter, logger)

	app := api.NewApplication(cfg, logger)
	app.SetRouter(api.NewRouter(handler))

	logger.Infow("Starting server", "port", cfg.Addr)
	if err := app.ServeHTTP(); err != nil {
		logger.Fatal(err)
	}
}
