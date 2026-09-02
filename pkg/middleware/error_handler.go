package middleware

import (
	"errors"
	"health-tracker-api/internal/model/web"
	"health-tracker-api/pkg/apperror"

	"github.com/gofiber/fiber/v3"
)

func GlobalErrorHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var notFoundErr *apperror.NotFoundError
	var validationErr *apperror.ValidationError
	var conflictErr *apperror.ConflictError
	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &notFoundErr):
		code = fiber.StatusNotFound
		message = notFoundErr.Message
	case errors.As(err, &validationErr):
		code = fiber.StatusBadRequest
		message = validationErr.Message
	case errors.As(err, &conflictErr):
		code = fiber.StatusConflict
		message = conflictErr.Message
	case errors.As(err, &fiberErr):
		code = fiberErr.Code
		message = fiberErr.Message
	default:
		message = err.Error()
	}

	return ctx.Status(code).JSON(web.WebResponse{
		Code: code,
		Status: false,
		Message: message,
	})
}