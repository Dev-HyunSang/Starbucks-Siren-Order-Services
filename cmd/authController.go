package cmd

import (
	"log"
	"time"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/dev-hyunsang/siren-order/models"
	"github.com/dongri/phonenumber"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "starbuck1971"

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	data := new(models.Users)
	users := new(models.Users)

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatal("Failed to DataBase Connection")
	}

	// 입력한 메일이 중복된 메일인지 확인함.
	db.Where("email = ?", data.Email).First(&users)
	if data.Email == users.Email {
		return c.Status(400).JSON(fiber.Map{
			"message": "중복되는 메일이 있습니다, 다시 확인 해 주세요.",
		})
	}

	// 입력한 닉네임이 중복된 닉네임인지 확인함.
	db.Where("nick_name =?", data.NickName).First(&users)
	if data.NickName == users.NickName {
		return c.Status(400).JSON(fiber.Map{
			"message": "중복되는 닉네임이 있습니다, 다시 확인 해 주세요.",
		})
	}

	number := phonenumber.Parse(data.PhoneNumber, "JP")
	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	uuid := uuid.NewV4()
	user := models.Users{
		UUID:        uuid,
		Name:        data.Name,
		NickName:    data.NickName,
		Birthday:    data.Birthday,
		PhoneNumber: number,
		Email:       data.Email,
		Password:    string(password),
		CreatedAt:   time.Now(),
	}

	db.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	req := new(RequestLogin)   // Request JSON
	users := new(models.Users) // Inset User Info

	if err := c.BodyParser(&req); err != nil {
		return err
	}
	log.Print(req)

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatalln("Failed to Connection DataBase")
	}

	// 사용자가 입력한 이메일을 바탕으로 정보 조회
	db.Where("email = ?", req.Email).Find(&users)
	log.Println(string(users.Email))
	log.Println(string(users.Password))

	// 만약 Email을 입력하지 않은 경우 오류 출력함.
	if req.Email == " " {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "메일 혹은 비밀번호가 올바르지 않습니다. 다시 확인 해 주세요.",
		})
	}

	// 입력한 비밀번호와 DB에 저장되어 있는 패스워드 대조 / 대조 실패시 오류 출력
	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(req.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "메일 혹은 비밀번호가 올바르지 않습니다. 다시 확인 해 주세요.",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    users.UUID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could Not Login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Login Success!",
	})
}

func EditUser(c *fiber.Ctx) error {
	req := new(models.Users)
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.Users

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatal("Failed to Connection DataaBase")

	}

	// JWT 저장되어 있는 JWT UUID를 통해서 정보를 가지고 옴.
	db.Where("uuid =?", claims.Issuer).Find(&user)

	//
	db.Model(&user).Updates(&models.Users{
		NickName: req.NickName,

		Email: req.Name,
	})

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "회원 정보가 성공적으로 수정되었습니다.",
	})
}

func Auth(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.Users

	db, err := database.ConnectionDataBase()
	if err != nil {
		log.Fatal("Failed to Connection DataaBase")

	}
	db.Where("uuid = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successed to Logout",
	})
}
