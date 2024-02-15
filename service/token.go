package service

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt      = "1j2h41lkj24lkjh5h"
	signinKey = "qwefq#12fasf3rf1q3"
	tokenTTL  = 12 * time.Hour
)

func (s *AuthService) generatePasswordHash(password string) (string, error) {

	hash := sha256.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil

}

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.SigningKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil

}
