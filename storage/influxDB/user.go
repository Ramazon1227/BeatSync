package influxdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	// "log"

	// "hash"
	"time"

	// cfg "github.com/Ramazon1227/BeatSync/config"
	"github.com/Ramazon1227/BeatSync/pkg/influxdb"
	"github.com/Ramazon1227/BeatSync/pkg/utils"
	"github.com/Ramazon1227/BeatSync/storage"
	"github.com/apache/arrow/go/v15/arrow"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	// "github.com/InfluxData/influxdb3-go/influxdb3"
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

	userID := uuid.NewV4().String()
	hashedPassword, err := utils.HashPassword(entity.Password)
	if err != nil {
		return nil, err
	}
	fmt.Println("Entity:", entity)

	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("user_info").
		SetTag("user_id", userID).
		SetTag("email", entity.Email).
		SetField("first_name", entity.FirstName).
		SetField("last_name", entity.LastName).
		SetField("password", hashedPassword).
		SetField("created_at", time.Now().UnixNano())

	if err := user.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return &models.PrimaryKey{
		Id: userID,
	}, nil

}

func (user *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	var userData = &models.User{}
	var count int
	// Execute query
	query := fmt.Sprintf(`SELECT *
                FROM "user_info"
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

func (user *UserRepoImpl) GetById(ctx context.Context, entity *models.PrimaryKey) (*models.User, error) {

	var userData = &models.User{}
	var count int
	// Execute query
	query := fmt.Sprintf(`SELECT *
                FROM "user_info"
                WHERE "user_id" = '%s'
				LIMIT 1`, entity.Id)

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
        _,ok:= value["user_id"]
		if ok {
			userData.ID = value["user_id"].(string)
		}
	
		_,ok = value["first_name"]
		if ok {
			userData.FirstName = value["first_name"].(string)
		}

		_,ok = value["last_name"]
		if ok {
			userData.LastName = value["last_name"].(string)
		}

		// userData.Email = value["email"].(string)
		_,ok = value["email"]
		if ok {
			userData.Email = value["email"].(string)
		}
		// userData.Password = value["password"].(string)
		// userData.Phone = value["phone"].(string)
		_,ok = value["phone"]
		if ok {
			userData.Phone = value["phone"].(string)
		}
		// userData.Gender = value["gender"].(string)
		_,ok = value["gender"]
		if ok {
			userData.Gender = value["gender"].(string)
		}
		// userData.Age = value["age"].(int)
		_,ok = value["age"]
		if ok {
			userData.Age = value["age"].(int64)
		}
		// userData.Height = value["height"].(float32)
		_,ok = value["height"]
		if ok {
			userData.Height = value["height"].(float64)
		}
		// userData.Weight = value["weight"].(float32)
		_,ok = value["weight"]
		if ok {
			userData.Weight = value["weight"].(float64)
		}
		_,ok = value["password"]
		if ok {
			userData.Password = value["password"].(string)
		}
		ts,ok := value["time"].(arrow.Timestamp)
		if ok {
			t := time.Unix(0, int64(ts)).UTC()
			userData.CreatedAt = &t
		}
	
	}
	if count == 0 {
		return nil, storage.ErrorNotFound
	}
	
	return userData, nil
}

func (user *UserRepoImpl) UpdateProfile(ctx context.Context, entity *models.UpdateProfileRequest) (pKey *models.PrimaryKey, err error) {

    // first get the user by ID we should get password and id 
	userData, err := user.GetById(ctx, &models.PrimaryKey{Id: entity.ID})
	if err != nil {
		log.Fatal("Error getting user by ID: ", err)
		return nil, err
	}
	
	fmt.Println("Entity:", entity)
	// delete the user by ID
	// err = user.Delete(ctx, userData.Email)
	// if err != nil {
	// 	log.Fatal("Error deleting user: ", err)
	// 	return nil, err
	// }

	
	
	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("user_info").
		SetTag("user_id", entity.ID).
		SetTag("email", userData.Email).
		SetField("first_name", entity.FirstName).
		SetField("last_name", entity.LastName).
		SetField("password", userData.Password).
		SetField("phone", entity.Phone).
		SetField("gender", entity.Gender).
		SetField("age", entity.Age).
		SetField("height", entity.Height).
		SetField("weight", entity.Weight).
		SetField("created_at", userData.CreatedAt.UnixNano()).
		SetField("updated_at", time.Now().UnixNano()).
		SetTimestamp(*userData.CreatedAt)


	if err := user.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
		panic(err)
	}

	return &models.PrimaryKey{
		Id: entity.ID,
	}, nil

}

func (user *UserRepoImpl) UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error {
	// Get the user by ID
	userData, err := user.GetById(ctx, &models.PrimaryKey{Id: userId})
	if err != nil {
		return err
	}

	// Check if the current password is correct
	if !utils.CheckPassword(userData.Password,newPassword) {
		return errors.New("current password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the password in the database
	options := influxdb3.WriteOptions{
		Database: "beatsync",
	}

	point := influxdb3.NewPointWithMeasurement("user_info").
		SetTag("user_id", userId).
		SetField("email", userData.Email).
		SetField("first_name", userData.FirstName).
		SetField("last_name", userData.LastName).
		SetField("password", hashedPassword).
		SetField("phone", userData.Phone).	
		SetField("gender", userData.Gender).
		SetField("age", userData.Age).
		SetField("height", userData.Height).
		SetField("weight", userData.Weight).
		SetField("created_at", userData.CreatedAt.UnixNano()).
		SetField("updated_at", time.Now().UnixNano()).
		SetTimestamp(*userData.CreatedAt)



		if err := user.db.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			panic(err)
		}

		return nil
}

func (user *UserRepoImpl) Delete(ctx context.Context, email string) (error) {

	err:= influxdb.DeleteUser(ctx, email)
	if err != nil {
		return err
	}
		
	return  nil
}