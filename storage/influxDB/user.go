package influxdb

import (
	"context"
	"fmt"
	"log"

	// "hash"
	"time"

	"github.com/Ramazon1227/BeatSync/config"
	// logger "github.com/Ramazon1227/BeatSync/pkg/logger"
	"github.com/Ramazon1227/BeatSync/pkg/utils"

	"github.com/Ramazon1227/BeatSync/models"
	// "github.com/influxdata/influxdb-client-go/api"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/satori/go.uuid"
)

type UserRepoImpl struct {
	db     *influxdb2.Client
	Org    string
	Bucket string
	// WriteAPI api.WriteAPI
	// QueryAPI api.QueryAPI
}

func NewUserRepoI(client influxdb2.Client) *UserRepoImpl {

	return &UserRepoImpl{
		db:     &client,
		Org:    config.Load().InfluxOrg,
		Bucket: config.Load().InfluxBucket,
		// // WriteAPI: client.WriteAPIBlocking(config.Load().InfluxOrg, config.Load().InfluxBucket),
		// // WriteAPI: (*client).WriteAPIBlocking(Config.InfluxOrg, Config.InfluxBucket),
		// 	WriteAPI: (*client).WriteAPI(Config.InfluxOrg, Config.InfluxBucket),
		// QueryAPI: (*client).QueryAPI(Config.InfluxOrg),
	}
}

func (user *UserRepoImpl) Create(ctx context.Context, entity *models.UserRegisterModel) (pKey *models.PrimaryKey, err error) {

	writeAPI := (*user.db).WriteAPIBlocking(user.Org, user.Bucket)

	// Generate a unique patient ID
	patientID := uuid.NewV4().String()
	hashedPassword, err := utils.HashPassword(entity.Password)
	if err != nil {
		return nil, err
	}
	// Create data point
	// data := fmt.Sprintf(
	// 	"user_info,patient_id=%s,first_name=%s,last_name=%s,password=%s,created_at =%d",
	// 	patientID, entity.FirstName, entity.LastName, hashedPassword, time.Now().UnixNano())

	tags := map[string]string{
		"email":      entity.Email,
		"patient_id": patientID,
	}
	fields := map[string]interface{}{
		"first_name": entity.FirstName,
		"last_name":  entity.LastName,
		"password":   hashedPassword,
		"created_at": time.Now().UnixNano(),
	}
	point := write.NewPoint("user_info", tags, fields, time.Now())

	// Write to InfluxDB
	err = writeAPI.WritePoint(context.Background(), point)

	if err != nil {
		log.Println("❌ Error writing user register data to InfluxDB:", err)
		return nil, err
	} else {
		fmt.Println("✅ user Data Saved to InfluxDB")
	}

	return &models.PrimaryKey{
		Id: patientID,
	}, nil

}

func (user *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {	

    var userData =&models.User{}
	queryAPI := (*user.db).QueryAPI(user.Org)
	query := fmt.Sprintf(`from(bucket: "%s") |> range(start: -10d) |> filter(fn: (r) => r._measurement == "user_info" and r._field == "email" and r.email == "%s")`, user.Bucket, email)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		log.Println("❌ Error querying InfluxDB storage level: getting user info by email ", err)
		return nil, err
	}
	defer result.Close()
   fmt.Println("Query result:", string(result.TablePosition()))
	if result.Next() {
		record := result.Record()
		userData := &models.User{
			ID:        int(record.ValueByKey("patient_id").(int)),
			FirstName: record.ValueByKey("first_name").(string),
			LastName:  record.ValueByKey("last_name").(string),
			Email:     email,
			Password:  record.ValueByKey("password").(string),
			CreatedAt: func(t time.Time) *time.Time { return &t }(time.Unix(record.Time().Unix(), 0)),
		}
		fmt.Println("I found user data in InfluxDB:", userData)
		return userData, nil
	}

	if result.Err() != nil {
		log.Println("❌ Error iterating over query results:", result.Err())
		return nil, result.Err()
	}

	return userData, nil
}