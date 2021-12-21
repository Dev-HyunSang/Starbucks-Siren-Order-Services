package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dev-hyunsang/siren-order/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/twinj/uuid"
)

type AccessDetails struct {
	AccessUUID string
	UserUUID   interface{}
}

func createJWTToken(userUUID uuid.UUID) (*TokenDetails, error) {
	td := TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = td.AccessUUID + "++" + userUUID.String()

	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	Claims := jwt.MapClaims{}
	Claims["authorized"] = true
	Claims["access_uuid"] = td.AccessUUID
	Claims["user_uuid"] = userUUID
	Claims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Fatalln("[ERROR] Failed to Create SignedString")
		return nil, err
	}

	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_uuid"] = userUUID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("REFRESH_SECRET"))
	if err != nil {
		log.Fatalln("[ERROR] Failed to Signed String")
	}
	return &td, nil
}

func CreateAuth(userUUID uuid.UUID, td *TokenDetails) error {
	client := database.Redis()
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUUID, userUUID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUUID, userUUID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.GetRespHeader("Authorization")
	log.Println(bearToken)
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return " "
}

func VerifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenVailed(c *fiber.Ctx) error {
	token, err := VerifyToken(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": token,
		})
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(c *fiber.Ctx) (*AccessDetails, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userUUID := claims["user_uuid"]

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserUUID:   userUUID,
		}, nil
	}
	return nil, err
}

func FetchAuth(autD *AccessDetails) (string, error) {
	client := database.Redis()
	userUUID, err := client.Get(autD.AccessUUID).Result()
	if err != nil {
		return "Failed to Fetch Auth at Redis Server", err
	}
	return userUUID, nil
}

func DeleteAuth(givenUUID string) (int64, error) {
	client := database.Redis()
	deleted, err := client.Del(givenUUID).Result()
	if err != nil {
		log.Fatalln("Fail to DeleteAuth at Redis Delete")
		return 0, err
	}
	return deleted, nil
}

func LogOut(c *fiber.Ctx) error {
	au, err := ExtractTokenMetadata(c)
	if err != nil {
		log.Fatal("Failed to ExtractTokenMetadata")
	}
	deleted, delError := DeleteAuth(au.AccessUUID)
	if delError != nil || deleted == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"data": "unauthorized",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func Refresh(c *fiber.Ctx) {
	mapToken := map[string]string{}
	if err := c.BodyParser(&mapToken); err != nil {
		log.Fatalln("Failed to Refresh JSON Body Parser")
		return
	}
	refreshToken := mapToken["refresh_token"]

	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		fmt.Println("the error: ", err)
		fiber.NewError(fiber.StatusUnauthorized, "Refresh token expired")
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"data":   err,
		})
		return
	}
	clamis, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUUID, ok := clamis["refresh_uuid"].(string)
		if !ok {
			c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status": fiber.StatusUnprocessableEntity,
				"data":   err,
			})
			return
		}

		deleted, delErr := DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": fiber.StatusUnauthorized,
				"data":   "Unauthorized",
			})
			return
		}
		ts, createErr := createJWTToken(uuid.NewV4())
		if createErr != nil {
			c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": fiber.StatusForbidden,
				"data":   createErr.Error(),
			})
			return
		}
		saveErr := CreateAuth(uuid.NewV4(), ts)
		if saveErr != nil {
			c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": fiber.StatusForbidden,
				"data":   saveErr.Error(),
			})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": fiber.StatusCreated,
			"data":   tokens,
		})
		return
	} else {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"data":   "refresh expired",
		})
		return
	}
}

func DeleteTokens(authD *AccessDetails) error {
	refreshUUID := fmt.Sprintf("%s++%d", authD.AccessUUID, authD.UserUUID)
	client := database.Redis()
	deleteAt, err := client.Del(authD.AccessUUID).Result()
	if err != nil {
		return err
	}
	deleteRt, err := client.Del(refreshUUID).Result()
	if err != nil {
		return err
	}
	if deleteAt != 1 || deleteRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
