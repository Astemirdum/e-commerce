package repo

import (
	"context"
	"github.com/Astemirdum/e-commerce/store-product/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	ProductRepository
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) (int64, error)
	FindOne(ctx context.Context, productId int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	CreateOrderLog(ctx context.Context, req *StockRequest) error
}

func NewProductRepository(db *gorm.DB) *Repository {
	return &Repository{ProductRepository: NewpDB(db)}
}

func PGConn(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.Product{}); err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.StockDecreaseLog{}); err != nil {
		return nil, err
	}
	return db, nil
}
