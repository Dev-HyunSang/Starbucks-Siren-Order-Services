package cmd

import (
	"log"
	"time"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/dev-hyunsang/siren-order/model"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

var JwtKey = []byte("starbucks")

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

	uuid := uuid.NewV4()

	user := &model.Users{
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

	token, err := createJWTToken(user.UUID)
	if err != nil {
		return err
	}

	cookie := new(fiber.Cookie)
	c.Response().Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	c.Response().Header.Set("Last-Modified", time.Now().String())
	c.Response().Header.Set("Pragma", "no-cache")
	c.Response().Header.Set("Expires", "-1")
	cookie.Name = "set-token"
	cookie.Value = token.AccessUUID

	c.Cookie(cookie)
	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"exp":    token,
		"isOk":   true,
	})
}

func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid login credentials")
	}

	db, err := database.ConnectionDataBase()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DataBase Connection ERROR")
	}

	users := new(model.Users)
	db.Where("email = ?", req.Email).Find(users)
	if users == nil {
		log.Print("Not Users Information at Null")
		return fiber.NewError(fiber.StatusInternalServerError, "Not Users Information at Null")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := createJWTToken(users.UUID)
	if err != nil {
		return err
	}

	saveErr := CreateAuth(users.UUID, token)
	if saveErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": saveErr,
		})
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tokens,
	})
}
