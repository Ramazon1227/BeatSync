package models

import "time"

type PPGData struct {
	TimeStamp *time.Time `json:"timestamp" binding:"required"`
	Value     float64    `json:"value" binding:"required"`
}

type SensorData struct {
	SensorDataID string     `json:"sensor_data_id"`               // Unique identifier for the sensor data
	UserID       string     `json:"user_id" binding:"required"`   // Unique identifier for the user
	DeviceID     string     `json:"device_id,omitempty"`          // Optional: Identifier for the specific device
	Time         *time.Time `json:"timestamp" binding:"required"`      // Timestamp of the reading
	Data         []PPGData `json:"data" binding:"required,gt=0"` // Heart rate value (e.g., beats per minute)
}

type SaveSensorDataResponse struct {
	SensorDataID string `json:"sensor_data_id"` // Unique identifier for the sensor data
	AnalysisID   string `json:"analysis_id"`    // Unique identifier for the analysis
}

type HRVAnalysisResult struct {
	AnalysisID   string  `json:"analysis_id"` // Unique identifier for the analysis
	UserID       string  `json:"user_id"`
	AnalysisTime string  `json:"analysis_time"` // Timestamp of when the analysis was performed
	SDNN         float64 `json:"sdnn"`
	RMSSD        float64 `json:"rmssd"`
	NN50         int64   `json:"nn50"`
	PNN50        float64 `json:"pnn50"`
	LF           float64 `json:"lf"`
	HF           float64 `json:"hf"`
	VLF          float64 `json:"vlf"`
	LFHF         float64 `json:"lf_hf_ratio"`
	SD1          float64 `json:"sd1"`
	SD2          float64 `json:"sd2"`
}

type UserHRVAnalysisRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type UserHRVAnalysisResponse struct {
	Count    int                  `json:"count"`
	Analysis []*HRVAnalysisResult `json:"analysis"`
}
