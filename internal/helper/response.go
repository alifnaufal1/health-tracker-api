package helper

import (
	"health-tracker-api/internal/model/web"

	"github.com/gofiber/fiber/v3"
)

func WriteResponseBody(c fiber.Ctx, response any) {
	c.RequestCtx().SetContentType("application/json")
	err := c.JSON(response)
	PanicIfError(err)
}

func ToWebResponse(c fiber.Ctx, data any) error {
	return c.Status(200).JSON(web.WebResponse{
		Code: 200,
		Status: true,
		Data: data,
	})
}