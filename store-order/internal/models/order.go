package models

type Order struct {
	Id        int64 `gorm:"primaryKey"`
	ProductId int64
	UserId    int64
	Failed    bool
}
