package repo

import (
	"context"

	"github.com/Astemirdum/e-commerce/store-auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(ctx context.Context, req *UserRequest) error
	Get(ctx context.Context, req *UserRequest) (*models.User, error)
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepository: NewpDB(db),
	}
}

func PGConn(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	return db, nil
}
