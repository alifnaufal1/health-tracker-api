package controller

import (
	"health-tracker-api/internal/helper"
	workoutDataWeb "health-tracker-api/internal/model/web/workout_data"
	"health-tracker-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type WorkoutDataController interface {
	Create(ctx fiber.Ctx) error
	// Update(ctx fiber.Ctx) error
	// Delete(ctx fiber.Ctx) error
	// FindByID(ctx fiber.Ctx) error
	// FindAll(ctx fiber.Ctx) error
}

type WorkoutDataControllerImpl struct {
	WorkoutDataService service.WorkoutDataService
}

func NewWorkoutDataController(workoutDataService service.WorkoutDataService, log *logrus.Logger) WorkoutDataController {
	return &WorkoutDataControllerImpl{
		WorkoutDataService: workoutDataService,
	}
}

func (c *WorkoutDataControllerImpl) Create(ctx fiber.Ctx) error {
	workoutDataCreateRequest := new(workoutDataWeb.WorkoutDataCreateRequest)
	err := ctx.Bind().Body(workoutDataCreateRequest)
	if err != nil {
		return err
	}
	
	createdWorkoutData, err := c.WorkoutDataService.Create(ctx, *workoutDataCreateRequest)
	if err != nil {
		return err
	}
	
	return helper.ToWebResponse(ctx, createdWorkoutData, "success create new workout data")
}