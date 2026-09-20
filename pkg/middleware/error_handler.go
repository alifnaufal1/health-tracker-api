package middleware

import (
	"errors"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func GlobalErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var notFoundErr *apperror.NotFoundError
	var validationErr *apperror.ValidationError
	var conflictErr *apperror.ConflictError
	var forbiddenErr *apperror.ForbiddenError
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
	case errors.As(err, &forbiddenErr):
		code = fiber.StatusForbidden
		message = forbiddenErr.Message
	case errors.As(err, &fiberErr):
		code = fiberErr.Code
		message = fiberErr.Message
	default:
		log.Errorf("unhandled error : %v",  err)
		message = "Internal Server Error"
	}

	return response.Error(c, code, message)
}