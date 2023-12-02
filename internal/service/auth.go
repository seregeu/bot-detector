package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int64
}

type Authorizer struct {
	config *Config
}

func NewAuthorizer() (*Authorizer, error) {
	configAuth, err := NewConfig()
	if err != nil {
		return nil, err
	}

	return &Authorizer{
		config: configAuth,
	}, nil
}

func (a *Authorizer) GenerateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.StandardClaims{
				IssuedAt: time.Now().Unix(),
				ExpiresAt: time.Now().Add(
					time.Hour * time.Duration(a.config.expireDuration),
				).Unix(),
			},
			userID,
		})

	signedToken, err := token.SignedString([]byte(a.config.signingKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (a *Authorizer) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.config.signingKey), nil
		})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
