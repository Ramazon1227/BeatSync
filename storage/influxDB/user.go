package influxdb

import (
	"context"
	// "fmt"

	"github.com/Ramazon1227/BeatSync/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type UserRepoImpl struct {
	db  *influxdb2.Client
}

func NewUserRepoI(client *influxdb2.Client) *UserRepoImpl {
	
	return &UserRepoImpl{db: client}
}

func (user *UserRepoImpl) Create (ctx context.Context , entity *models.UserCreateModel) (pKey *models.PrimaryKey, err error) {
    
	// writeAPI := (*user.db).WriteAPIBlocking(org, bucket)

	// // Create data point
	// data := fmt.Sprintf(
	// 	"hrv_analysis,patient_id=%s,age=%d,gender=%s,activity=%s sdnn=%.2f,rmssd=%.2f,heart_rate=%d,insight=\"%s\" %d",
	// 	patientID, age, gender, activity, sdnn, rmssd, heartRate, insight, time.Now().UnixNano())

	// // Write to InfluxDB
	// err := writeAPI.WriteRecord(context.Background(), data)
	// if err != nil {
	// 	log.Println("❌ Error writing HRV data to InfluxDB:", err)
	// } else {
	// 	fmt.Println("✅ HRV Data Saved to InfluxDB")
	// }

	return nil, nil
	
}