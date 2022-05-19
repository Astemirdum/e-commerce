package repo

import (
	"context"
	"github.com/Astemirdum/e-commerce/store-order/models"

	"gorm.io/gorm"
)

type pDB struct {
	DB *gorm.DB
}

func NewpDB(db *gorm.DB) *pDB {
	return &pDB{db}
}

func (db *pDB) CreateOrder(ctx context.Context, order *models.Order) error {
	return db.DB.Debug().WithContext(ctx).Create(order).Error
}

func (db *pDB) UpdateOrder(ctx context.Context, order *models.Order) error {
	return db.DB.Debug().WithContext(ctx).Save(order).Error
}
