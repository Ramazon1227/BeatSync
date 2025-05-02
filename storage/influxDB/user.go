package influxdb

import (
	"context"
	"fmt"

	// "log"

	// "hash"
	"time"

	// cfg "github.com/Ramazon1227/BeatSync/config"
	// logger "github.com/Ramazon1227/BeatSync/pkg/logger"
	"github.com/Ramazon1227/BeatSync/pkg/utils"
	"github.com/Ramazon1227/BeatSync/storage"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/Ramazon1227/BeatSync/models"
	uuid "github.com/satori/go.uuid"
)

type UserRepoImpl struct {
	db *influxdb3.Client
}

func NewUserRepoI(client *influxdb3.Client) *UserRepoImpl {

	return &UserRepoImpl{
		db: client,
	}
}

func (user *UserRepoImpl) Create(ctx context.Context, entity *models.UserRegisterModel) (pKey *models.PrimaryKey, err error) {

	patientID := uuid.NewV4().String()
	hashedPassword, err := utils.HashPassword(entity.Password)
	if err != nil {
		return nil, err
	}
	fmt.Println("Entity:", entity)

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("beatsync").
		SetTag("user_id", patientID).
		SetTag("email", entity.Email).
		SetField("first_name", entity.FirstName).
		SetField("last_name", entity.LastName).
		SetField("password", hashedPassword).
		SetField("created_at", time.Now().UnixNano())

	if err := user.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return &models.PrimaryKey{
		Id: patientID,
	}, nil

}

func (user *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	var userData = &models.User{}
	var count int
	// Execute query
	query := fmt.Sprintf(`SELECT *
                FROM "beatsync"
                WHERE "email" = '%s'
				LIMIT 1`, email)

	queryOptions := influxdb3.QueryOptions{
		Database: "beatsync",
	}
	iterator, err := user.db.QueryWithOptions(context.Background(), &queryOptions, query)

	if err != nil {
		panic(err)
	}
	
	for iterator.Next() {
		count++
		value := iterator.Value()

		userData.ID = value["user_id"].(string)
		userData.FirstName = value["first_name"].(string)
		userData.LastName = value["last_name"].(string)
		userData.Email = value["email"].(string)
		userData.Password = value["password"].(string)
		// userData.CreatedAt = time.Unix(0, value["created_at"].(int64))
	}
	if count == 0 {
		return nil, storage.ErrorNotFound
	}
	
	return userData, nil
}
