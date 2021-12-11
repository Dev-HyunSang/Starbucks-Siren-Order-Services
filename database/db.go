package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDataBase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "parkhs0625", "127.0.0.1:3306", "siren_order")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 구조체를 이용한 테이블 자동 생성
	db.AutoMigrate(&Users{})

	return db, err
}
