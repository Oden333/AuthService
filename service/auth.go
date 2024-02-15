package service

import (
	"Auth/models"
	"Auth/repository"
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	repo            repository.Authorization
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningKey      string
}

func NewAuthService(repo repository.Authorization, atTTL time.Duration, rtTTL time.Duration, skey string) *AuthService {
	return &AuthService{
		repo:            repo,
		AccessTokenTTL:  atTTL,
		RefreshTokenTTL: rtTTL,
		SigningKey:      skey,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input models.SignUpInput) error {
	passwordHash, err := s.generatePasswordHash(input.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     input.Name,
		Password: passwordHash,
		Email:    input.Email,
	}
	logrus.Debug("Got sign-up request. Debug Data:", user)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GetUser(ctx context.Context, input models.SignInInput) (Tokens, error) {
	passwordHash, err := s.generatePasswordHash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetUser(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return Tokens{}, err
		}
		return Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func (s *AuthService) GetUserByGUID(ctx context.Context, guid string) (Tokens, error) {

	user, err := s.repo.GetUserByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return Tokens{}, err
		}
		return Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	refreshTokenHash := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	logrus.Debugf("Got finding by refreshToken request, token: %s \n hashed token for mongo lookup: %s", refreshToken, refreshTokenHash)

	user, err := s.repo.GetByRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func (s *AuthService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	//Получаем токены
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(s.AccessTokenTTL).Unix(),
		Subject:   userId.Hex(),
	})
	res.AccessToken, err = accessToken.SignedString([]byte(s.SigningKey))
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = NewRefreshToken()
	if err != nil {
		return res, err
	}
	//	Refresh токен тип произвольный, формат передачи base64,
	//		хранится в базе исключительно в виде bcrypt хеша,
	//		должен быть защищен от изменения настороне клиента и попыток повторного использования.
	session := models.Session{
		//???
		RefreshToken: base64.StdEncoding.EncodeToString([]byte(res.RefreshToken)),
		ExpiresAt:    time.Now().Add(s.RefreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}
