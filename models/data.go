package models

import "time"

type PPGData struct {
	TimeStamp *time.Time `json:"timestamp" binding:"required"`
	Value     float64    `json:"value" binding:"required"`
}

type SensorData struct {
	UserID   string     `json:"user_id" binding:"required"`   // Unique identifier for the user
	DeviceID string     `json:"device_id,omitempty"`          // Optional: Identifier for the specific device
	Time     *time.Time `json:"time" binding:"required"`      // Timestamp of the reading
	Data     []*PPGData `json:"data" binding:"required,gt=0"` // Heart rate value (e.g., beats per minute)
}

type HRVAnalysisResult struct {
	UserID            string     `json:"user_id"`
	AnalysisTimestamp *time.Time `json:"analysis_time"` // Timestamp of when the analysis was performed
	SDNN              float64    `json:"sdnn,omitempty"`
	RMSSD             float64    `json:"rmssd,omitempty"`
	NN50              int        `json:"nn50,omitempty"`
	PNN50             float64    `json:"pnn50,omitempty"`
	LF                float64    `json:"lf,omitempty"`
	HF                float64    `json:"hf,omitempty"`
	VLF               float64    `json:"vlf,omitempty"`
	LFHF              float64    `json:"lf_hf_ratio,omitempty"`
	SD1               float64    `json:"sd1,omitempty"`
	SD2               float64    `json:"sd2,omitempty"`
}
