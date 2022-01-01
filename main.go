package main

import (
	"log"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/dev-hyunsang/siren-order/middleware"
	"github.com/dev-hyunsang/siren-order/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	middleware.Route(app)

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatalln(err)
	}
	// 구조체를 이용한 테이블 자동 생성
	db.AutoMigrate(
		&models.Users{},
		&models.Foods{},
		&models.Admin{},
	)

	if err := app.Listen(":4000"); err != nil {
		log.Fatalln("Fiber Listen Error!", err)
	}
}
