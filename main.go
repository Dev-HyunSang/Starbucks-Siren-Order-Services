package main

import (
	"log"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/dev-hyunsang/siren-order/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	middleware.Route(app)

	// DataBase Connection TEST at Ping
	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatalln("[ERROR] Failed DataBase Connection")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("[ERROR] Failed to SQL DataBase")
	}
	pingDB := sqlDB.Ping()
	log.Panicln("Ping DataBase", pingDB)

	// Auth / JWT Remote Dictionary Server TEST at Ping
	client := database.Redis()
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatalln("[ERROR] Failed to Redis Connect at Ping")
	}

	if err := app.Listen(":3000"); err != nil {
		log.Fatalln("Fiber Listen Error!", err)
	}
}
