package service

import (
	"errors"
	"time"

	"github.com/Astemirdum/e-commerce/store-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey string
	Issuer    string
	TokenTTL  time.Duration
}

type jwtClaims struct {
	jwt.StandardClaims
	UserId int64
	Email  string
}

func NewJwtWrapper(secretKey string, issuer string, tokenTTL time.Duration) *JwtWrapper {
	return &JwtWrapper{
		SecretKey: secretKey,
		Issuer:    issuer,
		TokenTTL:  tokenTTL, // minute
	}
}

func (jw *JwtWrapper) ParseToken(accessToken string) (*jwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(accessToken, &jwtClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("wrong signing method")
			}
			return []byte(jw.SecretKey), nil
		})
	if err != nil {
		return nil, err
	}
	myClaims, _ := jwtToken.Claims.(*jwtClaims)
	return myClaims, nil
}

func (jw *JwtWrapper) GenerateToken(user *models.User) (string, error) {
	claims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * jw.TokenTTL).UTC().Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    jw.Issuer,
		},
		UserId: user.Id,
		Email:  user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString([]byte(jw.SecretKey))
}
