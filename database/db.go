package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDataBase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "parkhs0625", "127.0.0.1:3306", "starbucks_siren_order")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return db, err
}
