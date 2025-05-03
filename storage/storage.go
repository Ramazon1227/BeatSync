package storage

import (
	"context"
	"errors"

	"github.com/Ramazon1227/BeatSync/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
}

var (
	ErrorNotFound error = errors.New("not found")
)

type UserRepoI interface {
	Create(ctx context.Context, entity *models.UserRegisterModel) (pKey *models.PrimaryKey, err error)
	GetByEmail(ctx context.Context, email string) (*models.User, error) // for login
	GetById(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error)
	UpdateProfile(ctx context.Context, entity *models.UpdateProfileRequest) (pkey *models.PrimaryKey,err error)
	// UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error
	Delete(ctx context.Context, email string) error
}

type AnalyzeRepoI interface {

	SaveSensorData(ctx context.Context, entity *models.SensorData) (pKey *models.PrimaryKey, err error)
}