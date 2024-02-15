package repository

import (
	"Auth/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) GetById(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User

	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}
