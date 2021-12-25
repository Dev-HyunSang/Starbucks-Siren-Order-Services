package main

import (
	"log"

	"github.com/dev-hyunsang/siren-order/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	middleware.Route(app)

	//database.ConnectionStorage()

	if err := app.Listen(":3000"); err != nil {
		log.Fatalln("Fiber Listen Error!", err)
	}
}
