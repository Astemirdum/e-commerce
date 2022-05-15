package repo

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/Astemirdum/e-commerce/store-auth/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	salt = "kjvbbe8392dsn"
)

type pDB struct {
	DB *gorm.DB
}

func NewpDB(db *gorm.DB) *pDB {
	return &pDB{db}
}

type UserRequest struct {
	Email    string
	Password string
}

var (
	ErrAlreadyExists = errors.New("user already exists")
)

func (db *pDB) Create(ctx context.Context, req *UserRequest) error {
	var user models.User

	if db.DB.WithContext(ctx).Where(&models.User{Email: req.Email}).First(nil).Error == nil {
		return ErrAlreadyExists
	}

	user.Email = req.Email
	user.Password = genPasswordHash(req.Password)
	return db.DB.Create(user).Error
}

func (db *pDB) Get(ctx context.Context, req *UserRequest) (*models.User, error) {
	var user *models.User

	if err := db.DB.WithContext(ctx).Where(&models.User{Email: req.Email}).First(user).Error; err != nil {
		return nil, errors.Errorf("user not authorized: %v", err)
	}
	return user, nil
}

func genPasswordHash(pass string) string {
	hash := sha1.New()
	hash.Write([]byte(pass))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
