package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollection = "users"
)
const timeout = 10 * time.Second

func NewMongoDB(uri, username, password string) (*mongo.Client, error) {
	conString := fmt.Sprintf(uri, username, password)
	opts := options.Client().ApplyURI(conString)
	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username, Password: password,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	/* 	dbs, err := client.ListDatabaseNames(ctx, bson.M{})
	   	if err != nil {
	   		return nil, err
	   	}
	   	fmt.Println(dbs) */
	client.Database("authService")
	logrus.Info("Successful on DB connection")
	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}
