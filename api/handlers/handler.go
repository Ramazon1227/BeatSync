package handlers

import (
	"time"

	"github.com/Ramazon1227/BeatSync/api/http"
	"github.com/Ramazon1227/BeatSync/config"
	"github.com/Ramazon1227/BeatSync/storage"

	"github.com/Ramazon1227/BeatSync/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     config.Config
	log     logger.LoggerI
	storage storage.StorageI
}

func NewHandler(cfg config.Config, log logger.LoggerI, svcs storage.StorageI) Handler {
	return Handler{
		cfg:     cfg,
		log:     log,
		storage: svcs,
	}
}

func (h *Handler) handleResponse(c *gin.Context, status http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 300:
		h.log.Info(
			"---Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	case code < 400:
		h.log.Warn(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	default:
		h.log.Error(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	}

	c.JSON(status.Code, http.Response{
		Status:      status.Status,
		Description: status.Description,
		Data:        data,
	})
}

func (h *Handler) getStartDateParam(c *gin.Context) (startDate string, err error) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0,0,-7).Format("2006-01-02"))
	if _, err := time.Parse("2006-01-02", startDateStr); err != nil {
		return "", err
	}
	return startDateStr, nil
}

func (h *Handler) getEndDateParam(c *gin.Context) (endDate string, err error) {
	endDateStr := c.DefaultQuery("end_date", time.Now().AddDate(0,0,1).Format("2006-01-02"))
	if _, err := time.Parse("2006-01-02", endDateStr); err != nil {
		return "", err
	}
	return endDateStr, nil
}
