package models

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Email    string
	Password string
}
