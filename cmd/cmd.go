package cmd

import (
	"log"
	"time"

	"github.com/dev-hyunsang/siren-order/config"
	"github.com/dev-hyunsang/siren-order/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	UUID      uuid.UUID
	Name      string     `json:"name"`
	NickName  string     `json:"nickname"`
	Birthday  *time.Time `json:"birthday" time_format:"2006-01-02" time_utc:"1"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func Home(c *fiber.Ctx) error {
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("Access-Control-Allow-Headers", "*")
	return c.JSON(fiber.Map{
		"title": "Welecome to Starbuck Siren Order",
	})
}

func Rregister(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if req.Name == "" || req.NickName == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid SigUp Credentials")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Bcrypt Generate From Password Error")
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Created UUID Error")
	}

	user := &database.Users{
		UUID:      uuid,
		Name:      req.Name,
		NickName:  req.NickName,
		Birthday:  req.Birthday,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	result, err := database.ConnectionDataBase()
	if err != nil {
		log.Print("[ERROR] DataBase Connection...")
		return fiber.NewError(fiber.StatusInternalServerError, "Connection DataBase Error")
	}

	result.Create(&user)

	token, exp, err := createJWTToken(*user)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"token": token,
		"exp":   exp,
		"user":  user,
		"uuid":  uuid,
	})
}

func createJWTToken(user database.Users) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_uuid"] = user.UUID
	claims["exp"] = exp
	t, err := token.SignedString([]byte(config.Config("secret")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, err
}
