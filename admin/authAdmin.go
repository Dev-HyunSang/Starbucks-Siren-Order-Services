package admin

import (
	"log"
	"time"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/dev-hyunsang/siren-order/models"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type adminLogin struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"`
	AdminName string    `json:"admin_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type adminRegister struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	OTPNumber string    `json:"otp_number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RandowOTPNubmer(c *)
func Register(c *fiber.Ctx) error {
	req := new(adminRegister)
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatalln("Failed to DataBase Connection")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 16)
	uuid := uuid.NewV4()
	admin := models.Admin{
		AdminUUID: uuid.String(),
		Name:      req.Name,
		ID:        req.ID,
		Password:  string(password),
		CreatedAt: time.Now(),
	}

	db.Create(&admin)

	return c.Status(200).JSON(admin)
}
