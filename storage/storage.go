package storage

import (
	// "context"
	// "github.com/Ramazon1227/BeatSync/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
}

type UserRepoI interface {
	// Add(ctx context.Context, entity *models.UserCreateModel) (pKey *models.PrimaryKey, err error)
	// UpdateProfile(ctx context.Context, entity *models.UpdateUserProfileModel) error
	// GetById(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error)
	// GetList(ctx context.Context, queryParam *models.QueryParam) (*models.GetUserListModel, error)
	// Delete(ctx context.Context, pKey *models.PrimaryKey) error
	// GetByEmail(ctx context.Context, email string) (*models.User, error)
	// UpdateUserProfile(ctx context.Context, userId string, req *models.UpdateProfileRequest) error
	// UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error
}
