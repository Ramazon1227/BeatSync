package main

import (
	// "context"

	"github.com/Ramazon1227/BeatSync/config"
	// "github.com/Ramazon1227/BeatSync/api"
	"github.com/Ramazon1227/BeatSync/pkg/logger"
	"github.com/gin-gonic/gin"
)



func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	// pgStore, err := postgres.NewPostgres(context.Background(), cfg)
	// if err != nil {
	// 	log.Panic("postgres.NewPostgres", logger.Error(err))
	// }
	// defer pgStore.CloseDB()

	// h := handlers.NewHandler(cfg, log, pgStore)

	// r := api.SetUpRouter(h, cfg)

	// r.Run(cfg.HTTPPort)
}