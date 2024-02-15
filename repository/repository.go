package repository

import (
	"Auth/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, email, password string) (models.User, error)
	SetSession(ctx context.Context, userId primitive.ObjectID, session models.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
}

type Users interface {
	GetById(ctx context.Context, id primitive.ObjectID) (models.User, error)
}

type Repository struct {
	Authorization
	Users
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db),
		Users:         NewUsersRepo(db),
	}
}
