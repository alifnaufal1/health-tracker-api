package service

import (
	"errors"
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/workout_data"
	"health-tracker-api/internal/repository"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type WorkoutDataService interface {
	Create(ctx fiber.Ctx, request web.WorkoutDataCreateRequest) (*web.WorkoutDataResponse, error)
	// Delete(ctx fiber.Ctx, workoutDataID string) error
	// FindByID(ctx fiber.Ctx, workoutDataID string) (*web.WorkoutDataResponse, error)
	// FindByAll(ctx fiber.Ctx) (*[]web.WorkoutDataResponse, error)
}

type WorkoutDataServiceImpl struct {
	workoutDataRepository repository.WorkoutDataRepository
	validate *validator.Validate
	log *logrus.Logger
}

func NewWorkoutDataService(workoutDataRepository repository.WorkoutDataRepository, validate *validator.Validate, log *logrus.Logger) WorkoutDataService {
	return &WorkoutDataServiceImpl{
		workoutDataRepository: workoutDataRepository,
		validate: validate,
		log: log,
	}
}

func (s *WorkoutDataServiceImpl) Create(ctx fiber.Ctx, request web.WorkoutDataCreateRequest) (*web.WorkoutDataResponse, error) {
	log := helper.LoggerWithRequestID(ctx, s.log)
	log.WithField("request", request).Info("Received create workout data request")
	
	err := s.validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for create workout data request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}
	
	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	totalSteps := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalSteps - request.WorkoutBatchData[0].TotalSteps
	totalDistance := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalDistance - request.WorkoutBatchData[0].TotalDistance
	totalCalories := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalCalories - request.WorkoutBatchData[0].TotalCalories
	heartRateAverage := 0
	for _, data := range request.WorkoutBatchData {
		heartRateAverage += data.HeartRate
	}
	heartRateAverage /= len(request.WorkoutBatchData)
	
	workoutData := &domain.WorkoutData{
		Base: domain.Base{ID: uuid.New()},
		WorkoutDataType: request.WorkoutDataType,
		TotalSteps: totalSteps,
		TotalDistance: totalDistance,
		TotalCalories: totalCalories,
		HeartRateAverage: heartRateAverage,
		Pace: request.Pace,
		Duration: request.Duration,
		WorkoutBatchData: request.WorkoutBatchData,
		DeviceId: request.DeviceID,
	}

	workoutData, err = s.workoutDataRepository.Save(ctx, tx, workoutData)
	if err != nil {
		log.Warn("Create workout data process failed at repository layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("workout_data_id", workoutData.ID).Info("workout data created successfully")

	return helper.ToWorkoutDataResponse(workoutData), nil
}


