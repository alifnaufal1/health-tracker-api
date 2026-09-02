package helper

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

func LoggerWithRequestID(c fiber.Ctx, log *logrus.Logger) *logrus.Entry {
	reqID, ok := c.Locals("requestid").(string)
	if !ok || reqID == "" {
		reqID = "unknown"
	}
	return log.WithField("request_id", reqID)
}