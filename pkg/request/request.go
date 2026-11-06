package request

import (
	"context"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/logging"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func BindAndValidate[T any](c fiber.Ctx, v *validator.Validate) (*T, error) {
	req := new(T)
	if err := c.Bind().Body(req); err != nil {
		return nil, &apperror.ValidationError{Message: "invalid request body"}
	}
	if err := v.Struct(req); err != nil {
		return nil, &apperror.ValidationError{Message: err.Error()}
	}
	return req, nil
}

func Context(c fiber.Ctx) context.Context {
	return logging.WithRequestID(c, requestid.FromContext(c))
}
