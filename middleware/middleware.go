package middleware

import (
	"github.com/dev-hyunsang/siren-order/cmd"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Route(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	member := app.Group("/api/member")
	member.Post("/register", cmd.Register)
	member.Post("/login", cmd.Login)
	member.Get("/auth", cmd.Auth)
	member.Post("/logout", cmd.Logout)

	home := app.Group("/api/home")
	home.Get("/comment", cmd.Comment)
}
