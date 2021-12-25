package cmd

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	"log"
	"time"
)

type Register struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

func GetClient() *auth.Client {
	opt := option.WithCredentialsFile("/Users/hyun.sang/dev/GitHub/Siren-order/config/siren-order-services-firebase-adminsdk-zb6ii-24d963d5a7.json")
	config := &firebase.Config{ProjectID: "siren-order-services"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	return client
}
func CreateUser(c *fiber.Ctx) error {
	req := new(Register)
	if err := c.BodyParser(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "[ERROR] Failed to Bucket Request at BodyParser",
			"data":    err,
			"time":    time.Now(),
		})
	}

	client := GetClient()
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		EmailVerified(false).
		PhoneNumber(req.PhoneNumber).
		Password(req.Password).
		DisplayName(req.Name).
		Disabled(false)
	u, err := client.CreateUser(context.Background(), params)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to CreateUser at Google FireBase",
			"time":    time.Now(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "Successfully Created User!",
		"data":    u,
		"time":    time.Now(),
	})
}
