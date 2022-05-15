package repo

import (
	"context"
	"errors"
	"github.com/Astemirdum/e-commerce/store-product/internal/models"
	"gorm.io/gorm"
)

type pDB struct {
	DB *gorm.DB
}

func NewpDB(db *gorm.DB) *pDB {
	return &pDB{db}
}

var (
	ErrAlreadyExists = errors.New("order log already exists")
)

func (db *pDB) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	if err := db.DB.WithContext(ctx).Create(product).Error; err != nil {
		return 0, err
	}
	return product.Id, nil
}

func (db *pDB) FindOne(ctx context.Context, productId int64) (*models.Product, error) {
	var product *models.Product
	if err := db.DB.WithContext(ctx).First(product, productId).Error; err != nil {
		return nil, err
	}
	return product, nil
}

type StockRequest struct {
	ProductId int64
	OrderId   int64
	Count     int64
}

func (db *pDB) UpdateProduct(ctx context.Context, product *models.Product) error {
	return db.DB.WithContext(ctx).Save(product).Error
}

func (db *pDB) CreateOrderLog(ctx context.Context, req *StockRequest) error {
	var log models.StockDecreaseLog

	if err := db.DB.WithContext(ctx).Where(&models.StockDecreaseLog{OrderId: req.OrderId}).
		First(&log).Error; err == nil {
		return ErrAlreadyExists
	}

	log.OrderId = req.OrderId
	log.ProductRefer = req.ProductId
	log.Count = req.Count

	return db.DB.WithContext(ctx).Create(&log).Error
}

func (db *pDB) DecreaseStock(ctx context.Context, req *StockRequest) error {
	product, err := db.FindOne(ctx, req.ProductId)
	if err != nil {
		return err
	}

	if product.Stock >= req.Count {
		product.Stock -= req.Count
	}
	if err = db.UpdateProduct(ctx, product); err != nil {
		return err
	}

	if err = db.CreateOrderLog(ctx, req); err != nil {
		return err
	}

	return nil
}
