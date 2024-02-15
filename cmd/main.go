package main

import (
	server "Auth"
	"Auth/configs"
	"Auth/handler"
	"Auth/repository"
	"Auth/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	//Инициализируем логгер
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	//Инициализируем конфиги
	cfg, err := configs.InitConfig()
	if err != nil {
		logrus.Fatalf("unable to get config (%s)", err.Error())
	}

	//Подключаем БД
	mongoClient, err := repository.NewMongoDB(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logrus.Fatalf("Failed to initialize db (%s)", err.Error())
	}
	db := mongoClient.Database(cfg.Mongo.Name)
	//Зависимости чистой архитектуры
	repos := repository.NewRepository(db)
	services := service.NewService(service.Deps{
		Repo:            repos,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
		SigningKey:      cfg.Auth.JWT.SigningKey,
	})
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {

		if err := srv.Run(cfg.HTTP.Port, handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("Auth Service Started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logrus.Error(err.Error())
	}
}
