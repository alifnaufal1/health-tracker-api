package workout_data

import (
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type WorkoutDataHandler interface {
	Create(c fiber.Ctx) error
	GetByDeviceID(c fiber.Ctx) error
}

type WorkoutDataHandlerImpl struct {
	WorkoutDataService WorkoutDataService
}

func NewWorkoutDataHandler(workoutDataService WorkoutDataService) WorkoutDataHandler {
	return &WorkoutDataHandlerImpl{
		WorkoutDataService: workoutDataService,
	}
}

func (h *WorkoutDataHandlerImpl) Create(c fiber.Ctx) error {
	workoutDataCreateRequest := new(WorkoutDataCreateRequest)
	err := c.Bind().Body(workoutDataCreateRequest)
	if err != nil {
		return err
	}

	workoutDataCreateRequest.DeviceID = c.Params("deviceId")
	
	createdWorkoutData, err := h.WorkoutDataService.Create(c, *workoutDataCreateRequest)
	if err != nil {
		return err
	}
	
	return response.Success(c, createdWorkoutData, "success create new workout data")
}

func (h *WorkoutDataHandlerImpl) GetByDeviceID(c fiber.Ctx) error {	
	workoutData, err := h.WorkoutDataService.GetAll(c, c.Params("deviceId"))
	if err != nil {
		return err
	}
	
	return response.Success(c, workoutData, "success get all workout data")
}

