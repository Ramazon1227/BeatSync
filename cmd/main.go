package main

import (
	// "context"

	"context"

	"github.com/Ramazon1227/BeatSync/api"
	"github.com/Ramazon1227/BeatSync/api/handlers"
	"github.com/Ramazon1227/BeatSync/config"
	"github.com/Ramazon1227/BeatSync/pkg/logger"
	influxdb "github.com/Ramazon1227/BeatSync/storage/influxDB"
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

	influxdb, err := influxdb.NewInfluxDB(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer influxdb.CloseDB()

	h := handlers.NewHandler(cfg, log, influxdb)

	r := api.SetUpRouter(h, cfg)

	r.Run(cfg.HTTPPort)
}