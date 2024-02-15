package repository

import (
	"Auth/models"
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Collection
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{
		db: db.Collection(usersCollection),
	}
}

func (r *AuthMongo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if IsDuplicate(err) {
		return ErrorUserAlreadyExists
	}
	logrus.Debug("User created: ", user)
	return err
}

func (r *AuthMongo) GetUser(ctx context.Context, email, password string) (models.User, error) {
	var user models.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *AuthMongo) SetSession(ctx context.Context, userId primitive.ObjectID, session models.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"session": session}})

	return err
}

func (r *AuthMongo) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	var user models.User

	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}
