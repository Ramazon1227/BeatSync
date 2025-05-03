package influxdb

import (
	"context"
	"fmt"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/Ramazon1227/BeatSync/models"
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

func (data *AnalyzeRepoImpl) SaveSensorData(ctx context.Context, entity *models.SensorData) ( err error) {

	
	fmt.Println("Entity:", entity)

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("sensor_data").
		SetTag("user_id", entity.UserID).
		SetTag("device_id", entity.DeviceID).
		SetField("date", entity.Time).
		SetField("data", entity.Data).
		SetField("created_at", time.Now().UnixNano())

	if err := data.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return  nil

}

func (data *AnalyzeRepoImpl) SaveAnalysis(ctx context.Context, entity *models.SensorData) (pkey *models.PrimaryKey, err error) {
   
	analysisID := uuid.NewV4().String()
	
	fmt.Println("Entity:", entity)

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("analysis_data").
		SetTag("analysis_id", analysisID).
		SetTag("user_id", entity.UserID).
		SetField("date", entity.Time).
		SetField("data", entity.Data).
		SetField("created_at", time.Now().UnixNano())

	if err := data.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return  &models.PrimaryKey{Id: analysisID}, nil

}



