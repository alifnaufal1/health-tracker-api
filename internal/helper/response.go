package helper

import "github.com/gofiber/fiber/v3"

func WriteResponseBody(c fiber.Ctx, response any) {
	c.RequestCtx().SetContentType("application/json")
	err := c.JSON(response)
	PanicIfError(err)
}