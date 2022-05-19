package jwtoken

import (
	"errors"
	"time"

	"github.com/Astemirdum/e-commerce/store-auth/models"
	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey string
	Issuer    string
	TokenTTL  time.Duration
}

type JwtClaims struct {
	jwt.StandardClaims
	UserId int64
	Email  string
}

//go:generate mockgen -source=jwt.go -destination=mocks/mock.go

type JwtToken interface {
	ParseToken(accessToken string) (*JwtClaims, error)
	GenerateToken(user *models.User) (string, error)
}

func NewJwtWrapper(secretKey string, issuer string, tokenTTL time.Duration) *JwtWrapper {
	return &JwtWrapper{
		SecretKey: secretKey,
		Issuer:    issuer,
		TokenTTL:  tokenTTL, // minute
	}
}

func (jw *JwtWrapper) ParseToken(accessToken string) (*JwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(accessToken, &JwtClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("wrong signing method")
			}
			return []byte(jw.SecretKey), nil
		})
	if err != nil {
		return nil, err
	}
	myClaims, _ := jwtToken.Claims.(*JwtClaims)
	return myClaims, nil
}

func (jw *JwtWrapper) GenerateToken(user *models.User) (string, error) {
	claims := JwtClaims{
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
