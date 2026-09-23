package response

import "github.com/gofiber/fiber/v3"

type WebResponse struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type WebResponseError struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Error string `json:"error"`
}

func Success(c fiber.Ctx, data any, message string, ) error {
	return c.Status(200).JSON(WebResponse{
		Code:    200,
		Status:  true,
		Data:    data,
		Message: message,
	})
}

func Error(c fiber.Ctx, code int, errorMsg string) error {
	return c.Status(code).JSON(WebResponseError{
		Code:    code,
		Status:  false,
		Error: errorMsg,
	})
}

func ToResponses[A, B any](data *[]A, toResponse func(*A) *B) *[]B {
	var dataResponses []B
	for i := range *data {
		dataResponses = append(dataResponses, *toResponse(&(*data)[i]))
	}
	return &dataResponses
}