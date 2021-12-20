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

type Users struct {
	UserUUID  uuid.UUID
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
	req := new(Users)
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

	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"isOk":   true,
	})
}

func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	users := new(Users)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid login credentials")
	}

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatalln("Failed to DataBase Connection")
	}

	result := db.Where("email=?", req.Email).Find(&users)
	log.Println(result)
	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(req.Email))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "메일과 비밀번호를 확인 해 주세요.",
			"time":    time.Now(),
		})
	} else {
		ts, err := createJWTToken(uuid.NewV4())
		if err != nil {
			log.Fatal("Failed to Create JWT Token")
		}

		saveErr := CreateAuth(users.UserUUID, ts)
		if saveErr != nil {
			log.Fatal("Failed to Create Auth")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Failed to Create Auth",
				"data":    err,
			})
		}
		token := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
			"time":  time.Now(),
		})

	}
	return nil
}
