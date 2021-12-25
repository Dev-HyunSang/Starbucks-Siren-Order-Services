package cmd

import (
	"time"

	"github.com/NaySoftware/go-fcm"
	"github.com/dev-hyunsang/siren-order/config"
	"github.com/gofiber/fiber/v2"
)

type ReqeustFCM struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	SubMessage string `json:"sub_message"`
	Topic      string `json:"topic"`
}

func SendFCM(c *fiber.Ctx) error {
	req := new(ReqeustFCM)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "[ERROR] BodyParser SendFCM Reqeust Struct",
			"error":   err,
			"time":    time.Now(),
		})
	}
	data := map[string]string{
		"title": req.Title,
		"msg":   req.Message,
		"sum":   req.SubMessage,
	}
	serverKey := config.Config("FCM_SERVER_KEY")
	fcmClient := fcm.NewFcmClient(serverKey)
	fcmClient.NewFcmMsgTo(req.Topic, data)

	status, err := fcmClient.Send()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  status,
			"message": "[ERROR] Failed to FCM Message Send",
			"error":   err,
			"time":    time.Now(),
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  status,
		"message": "[Success] Successfully sent a message through FCM",
		"time":    time.Now(),
	})
}
