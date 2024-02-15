package configs

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP  HTTPConfig
		Mongo MongoConfig
		Auth  AuthConfig
	}

	MongoConfig struct {
		URI      string `mapstructure:"uri"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"pswd"`
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

func InitConfig() (*Config, error) {

	return &Config{
		HTTP: HTTPConfig{
			Host: "localhost",
			Port: "8080",
		},
		Mongo: MongoConfig{
			URI:      "mongodb://%s:%s@localhost:27019/",
			User:     "root",
			Password: "qwerty",
			Name:     "authService",
		},
		Auth: AuthConfig{
			JWT: JWTConfig{
				AccessTokenTTL:  time.Minute * 1000,
				RefreshTokenTTL: time.Hour * 720,
				SigningKey:      "40d084e7e4df17c35aaf91c5d1b5a6384dc1d28b2f8a6d50161b7d37a25bc31a",
			},
		},
	}, nil
	/*
			jwt:
		  signingKey: "40d084e7e4df17c35aaf91c5d1b5a6384dc1d28b2f8a6d50161b7d37a25bc31a"

		auth:
		  accessTokenTTL: 10m
		  refreshTokenTTL: 720h

		mongo:
		  uri: mongodb://%s:%s@localhost:27019/
		  user: root
		  pswd: qwerty
		  databaseName: authService
		  port: 27019

			var cfg Config
			return
			viper.AddConfigPath("configs")
			viper.SetConfigName("config")
			err := viper.ReadInConfig()
			if err != nil {
				return nil, errors.New("unable to read config")
			}
			if err = cfg.getCfg(); err != nil {
				return nil, errors.New("unable to get config data")
			}
			return &cfg, nil */
}

func (cfg *Config) getCfg() error {

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt.signingKey", &cfg.Auth.JWT.SigningKey); err != nil {
		return err
	}
	return nil
}
