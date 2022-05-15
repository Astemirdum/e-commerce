package models

type Product struct {
	Id                int64 `gorm:"primaryKey"`
	Sku               string
	Stock             int64
	Price             int64
	StockDecreaseLogs StockDecreaseLog `gorm:"foreignKey:ProductRefer"`
}

type StockDecreaseLog struct {
	Id           int64 `gorm:"primaryKey"`
	OrderId      int64
	Count        int64
	ProductRefer int64
}
