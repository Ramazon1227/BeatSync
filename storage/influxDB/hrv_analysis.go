package influxdb

import (
	"context"
	"fmt"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/Ramazon1227/BeatSync/models"
	"github.com/Ramazon1227/BeatSync/pkg/hrv"
	"github.com/Ramazon1227/BeatSync/storage"
	uuid "github.com/satori/go.uuid"
)

type AnalyzeRepoImpl struct {
	db *influxdb3.Client
}

func NewAnalyzeRepo(client *influxdb3.Client) *AnalyzeRepoImpl {

	return &AnalyzeRepoImpl{
		db: client,
	}
}

func (data *AnalyzeRepoImpl) SaveSensorData(ctx context.Context, entity *models.SensorData) (pKey *models.PrimaryKey, err error) {

	fmt.Println("Entity:", entity)
	sensorDataID := uuid.NewV4().String()

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("sensor_data").
		SetTag("user_id", entity.UserID).
		SetTag("device_id", entity.DeviceID).
		SetTag("sensor_data_id", sensorDataID).
		SetField("date", entity.Time).
		SetField("data", entity.Data).
		SetField("bpm", entity.BPM).
		SetField("created_at", time.Now().UnixNano())

	if err := data.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return &models.PrimaryKey{Id: sensorDataID}, nil

}

func (data *AnalyzeRepoImpl) SaveAnalysis(ctx context.Context, entity *models.SensorData) (pkey *models.PrimaryKey, err error) {

	analysisID := uuid.NewV4().String()
	fmt.Println("Entity:", entity)
	sensorDataID := uuid.NewV4().String()
	rr := hrv.ExtractRR(entity.Data, 300, 0.7) // minPeakDistanceMs=300, minPeakHeight=0.5
   fmt.Println("RR Intervals:", rr)
	sdnn := hrv.CalculateSDNN(rr)
	rmssd := hrv.CalculateRMSSD(rr)
	nn50 := hrv.CalculateNN50(rr)
	pnn50 := hrv.CalculatePNN50(rr)
	sd1, sd2 := hrv.CalculateSD1SD2(rr)
	lf, hf, vlf, lfHfRatio := hrv.CalculateFrequencyDomain(rr) 

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("analysis_data").
		SetTag("analysis_id", analysisID).
		SetTag("user_id", entity.UserID).
		SetTag("device_id", entity.DeviceID).
		SetTag("sensor_data_id", sensorDataID).
		SetField("analysis_time", time.Now()).
		SetField("bpm", entity.BPM).
		SetField("sdnn", sdnn).
		SetField("rmssd", rmssd).
		SetField("nn50", nn50).
		SetField("pnn50", pnn50).
		SetField("lf", lf).
		SetField("hf", hf).
		SetField("vlf", vlf).
		SetField("lf_hf_ratio", lfHfRatio).
		SetField("sd1", sd1).
		SetField("sd2", sd2).
		SetField("created_at", time.Now().UnixNano())
	if err := data.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return &models.PrimaryKey{Id: analysisID}, nil

}

func (data *AnalyzeRepoImpl) GetAnalysisByID(ctx context.Context,pkey *models.PrimaryKey) (*models.HRVAnalysisResult, error) {
	query := fmt.Sprintf(`SELECT * FROM "analysis_data" WHERE "analysis_id" = '%s' LIMIT 1`, pkey.Id)

	queryOptions := influxdb3.QueryOptions{
		Database: "beatsync",
	}

	iterator, err := data.db.QueryWithOptions(ctx, &queryOptions, query)
	if err != nil {
		return nil, err
	}

	if iterator.Next() {
		value := iterator.Value()
		analysis := &models.HRVAnalysisResult{
			AnalysisID:   value["sensor_data_id"].(string),
			UserID:       value["user_id"].(string),
			SDNN:         value["sdnn"].(float64),
			RMSSD:        value["rmssd"].(float64),
			NN50:         value["nn50"].(int64),
			PNN50:        value["pnn50"].(float64),
			SD1:          value["sd1"].(float64),
			SD2:          value["sd2"].(float64),
			LF:           value["lf"].(float64),
			HF:           value["hf"].(float64),
			VLF:          value["vlf"].(float64),
			LFHF:         value["lf_hf_ratio"].(float64),
			BPM: 		  value["bpm"].(int64),
			AnalysisTime: value["analysis_time"].(string),
		}
		return analysis, nil
	}

	return nil, storage.ErrorNotFound
}

func (data *AnalyzeRepoImpl) GetUserAnalysis(ctx context.Context, userID, startDate, endDate string) (*models.UserHRVAnalysisResponse, error) {
	query := fmt.Sprintf(`SELECT * FROM "analysis_data" WHERE "user_id" = '%s'`, userID)

	// Add date filters if provided
	if startDate != "" {
		query += fmt.Sprintf(` AND "analysis_time" >= '%s'`, startDate)
	}
	if endDate != "" {
		query += fmt.Sprintf(` AND "analysis_time" <= '%s'`, endDate)
	}

	queryOptions := influxdb3.QueryOptions{
		Database: "beatsync",
	}
    fmt.Println("Query:", query)
	iterator, err := data.db.QueryWithOptions(ctx, &queryOptions, query)
	if err != nil {
		return nil, err
	}
    
	var analysisList []*models.HRVAnalysisResult
	var count int
	for iterator.Next() {
		count++
		value := iterator.Value()
		analysis := models.HRVAnalysisResult{
			AnalysisID:   value["analysis_id"].(string),
			UserID:       value["user_id"].(string),
			BPM:          value["bpm"].(int64),
			SDNN:         value["sdnn"].(float64),
			RMSSD:        value["rmssd"].(float64),
			NN50:         value["nn50"].(int64),
			PNN50:        value["pnn50"].(float64),
			SD1:          value["sd1"].(float64),
			SD2:          value["sd2"].(float64),
			LF:           value["lf"].(float64),
			HF:           value["hf"].(float64),
			VLF:          value["vlf"].(float64),
			LFHF:         value["lf_hf_ratio"].(float64),
			AnalysisTime: value["analysis_time"].(string),
		}
		analysisList = append(analysisList, &analysis)
	}

	return &models.UserHRVAnalysisResponse{
		Count:    count,
		Analysis: analysisList,
	}, nil
}
