package main

import (
	"github.com/dev-hyunsang/siren-order/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	middleware.Route(app)

	if err := app.Listen(":3000"); err != nil {
		panic("Fiber Listen Error!")
	}
}
