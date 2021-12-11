package middleware

import (
	"github.com/dev-hyunsang/siren-order/cmd"
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Post("/", cmd.Home)
	app.Post("/register", cmd.Rregister)
	app.Post("/login", cmd.Login)
}
