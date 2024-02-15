package service

import (
	"Auth/models"
	"Auth/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo            repository.Users
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, AccessTokenTTL time.Duration, RefreshTokenTTL time.Duration) *UserService {

	return &UserService{
		repo:            repo,
		AccessTokenTTL:  AccessTokenTTL,
		RefreshTokenTTL: RefreshTokenTTL,
	}
}
func (s *UserService) GetUserById(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	return s.repo.GetById(ctx, id)
}
