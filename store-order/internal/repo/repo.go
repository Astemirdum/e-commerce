package repo

import (
	"context"

	"github.com/Astemirdum/e-commerce/store-order/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	OrderRepository
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	UpdateOrder(ctx context.Context, order *models.Order) error
}

func NewOrderRepository(db *gorm.DB) *Repository {
	return &Repository{OrderRepository: NewpDB(db)}
}

func PGConn(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.Order{}); err != nil {
		return nil, err
	}
	return db, nil
}
