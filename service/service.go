package service

import (
	"Auth/models"
	"Auth/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.SignUpInput) error
	GetUser(ctx context.Context, input models.SignInInput) (Tokens, error)
	GetUserByGUID(ctx context.Context, guid string) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	ParseToken(accessToken string) (string, error)
}
type Users interface {
	GetUserById(ctx context.Context, id primitive.ObjectID) (models.User, error)
}

type Service struct {
	Authorization
	Users
}

type Deps struct {
	Repo            *repository.Repository
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningKey      string
}

func NewService(deps Deps) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repo, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.SigningKey),
		Users:         NewUsersService(deps.Repo, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
