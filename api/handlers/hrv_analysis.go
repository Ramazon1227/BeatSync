package handlers

import (
	"context"

	httpapi "github.com/Ramazon1227/BeatSync/api/http"
	"github.com/Ramazon1227/BeatSync/models"
	"github.com/Ramazon1227/BeatSync/storage"

	"github.com/Ramazon1227/BeatSync/pkg/utils"
	"github.com/gin-gonic/gin"
)

// SaveSensorData godoc
// @ID save-sensor-data
// @Router /v1/sensor-data [POST]
// @Summary Save Sensor Data
// @Description Save sensor data for analysis
// @Tags data
// @Accept json
// @Produce json
// @Param sensor_data body models.SensorData true "Sensor data"
// @Success 201 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) SaveSensorData(c *gin.Context) {
    var sensorData models.SensorData

    err := c.ShouldBindJSON(&sensorData)
    if err != nil {
        h.handleResponse(c, httpapi.BadRequest, err.Error())
        return
    }

    respData,err := h.storage.Analyse().SaveSensorData(context.Background(), &sensorData)
    if err != nil {
        h.handleResponse(c, httpapi.InternalServerError, err)
        return
    }
	if respData == nil {
		h.handleResponse(c, httpapi.InternalServerError, "Failed to save sensor data")
		return
	}

	// Save analysis data	
	analysisResp, err := h.storage.Analyse().SaveAnalysis(context.Background(), &sensorData)
	if err != nil {
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}
	if analysisResp == nil {
		h.handleResponse(c, httpapi.InternalServerError, "Failed to save analysis data")
		return
	}


    h.handleResponse(c, httpapi.Created, models.SaveSensorDataResponse{
		SensorDataID: respData.Id,	
		AnalysisID:   analysisResp.Id,
	})
}

// GetAnalysisByID godoc
// @ID get-analysis-by-id
// @Router /v1/data/analysis/{analysis_id} [GET]
// @Summary Get Analysis By ID
// @Description Retrieve analysis data by its ID
// @Tags data
// @Accept json
// @Produce json
// @Param analysis_id path string true "Analysis ID"
// @Success 200 {object} models.Analysis
// @Failure 400 {object} httpapi.Response
// @Failure 404 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) GetAnalysisByID(c *gin.Context) {
    analysisID := c.Param("analysis_id")
    if analysisID == "" {
        h.handleResponse(c, httpapi.BadRequest, "Analysis ID is required")
        return
    }

    analysis, err := h.storage.Analyse().GetAnalysisByID(context.Background(), &models.PrimaryKey{Id: analysisID})
    if err != nil {
        if err == storage.ErrorNotFound {
            h.handleResponse(c, httpapi.NoContent, "Analysis not found")
            return
        }
        h.handleResponse(c, httpapi.InternalServerError, err)
        return
    }

    h.handleResponse(c, httpapi.OK, analysis)
}


// GetUserAnalysis godoc
// @ID get-user-analysis
// @Router /v1/data/user/{user_id}/analysis [GET]
// @Summary Get User Analysis
// @Description Retrieve all analysis data for a specific user within a date range
// @Tags data
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param start_date query string false "Start date in YYYY-MM-DD format"
// @Param end_date query string false "End date in YYYY-MM-DD format"
// @Success 200 {array} models.Analysis
// @Failure 400 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
func (h *Handler) GetUserAnalysis(c *gin.Context) {
    userID := c.Param("user_id")
    if userID == "" {
        h.handleResponse(c, httpapi.BadRequest, "User ID is required")
        return
    }

    // Parse query parameters for date range
    startDate,err := h.getStartDateParam(c)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err.Error())
		return
	}
	endDate,err := h.getEndDateParam(c)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err.Error())
		return
	}

    // Validate date format if provided
    if startDate != "" && !utils.IsValidDate(startDate) {
        h.handleResponse(c, httpapi.BadRequest, "Invalid start_date format. Use YYYY-MM-DD")
        return
    }
    if endDate != "" && !utils.IsValidDate(endDate) {
        h.handleResponse(c, httpapi.BadRequest, "Invalid end_date format. Use YYYY-MM-DD")
        return
    }

    // Fetch analysis data with optional date filters
    analysisResp, err := h.storage.Analyse().GetUserAnalysis(context.Background(), userID, startDate, endDate)
    if err != nil {
        h.handleResponse(c, httpapi.InternalServerError, err)
        return
    }

    h.handleResponse(c, httpapi.OK, analysisResp)
}