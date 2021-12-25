package middleware

import (
	"github.com/dev-hyunsang/siren-order/cmd"
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	member := app.Group("/member")
	member.Post("/register", cmd.Rregister)
	member.Post("/login", cmd.Login)
}
