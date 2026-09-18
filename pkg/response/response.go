package response

import "github.com/gofiber/fiber/v3"

type WebResponse struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Success(c fiber.Ctx, data any, message string, ) error {
	return c.Status(200).JSON(WebResponse{
		Code:    200,
		Status:  true,
		Data:    data,
		Message: message,
	})
}

func Error(c fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(WebResponse{
		Code:    code,
		Status:  false,
		Data:    nil,
		Message: message,
	})
}

func ToResponses[A any, B any](data *[]A, toResponse func(*A) *B) *[]B {
	var dataResponses []B
	for i := range *data {
		dataResponses = append(dataResponses, *toResponse(&(*data)[i]))
	}
	return &dataResponses
}