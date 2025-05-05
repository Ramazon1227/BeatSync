package storage

import (
	"context"
	"errors"

	"github.com/Ramazon1227/BeatSync/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
	Analyse() AnalyzeRepoI
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
	SaveAnalysis(ctx context.Context, entity *models.SensorData) (pkey *models.PrimaryKey, err error)
	GetAnalysisByID(ctx context.Context, pKey *models.PrimaryKey) (*models.HRVAnalysisResult, error)
	GetUserAnalysis(ctx context.Context, userID, startDate, endDate string) (*models.UserHRVAnalysisResponse, error)
}