package cmd

import (
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
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	Claims := jwt.MapClaims{}
	Claims["authorized"] = true
	Claims["access_uuid"] = td.AccessUUID
	Claims["user_uuid"] = userUUID
	Claims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Fatalln("[ERROR] Failed to Create SignedString")
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
	bearToken := c.Request().Header.Header()
	log.Println(bearToken)
	strArr := strings.Split(string(bearToken), " ")
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

func ExtractTokenMetadata(r *fiber.Ctx) (*AccessDetails, error) {
	token, err := VerifyToken(r)
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
