package workout_data

import (
	"fmt"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/model"
	"health-tracker-api/pkg/response"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceChecker interface {
	IsOwnedByUser(c fiber.Ctx, tx *gorm.DB, deviceID string, userID uuid.UUID) (bool, error)
}

type WorkoutDataService interface {
	Create(c fiber.Ctx, request WorkoutDataCreateRequest) (*WorkoutDataResponse, error)
	GetByDeviceID(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error)
}

type WorkoutDataServiceImpl struct {
	workoutDataRepository WorkoutDataRepository
	validate *validator.Validate
	log *logrus.Entry
	deviceChecker DeviceChecker 
}

func NewWorkoutDataService(workoutDataRepository WorkoutDataRepository, validate *validator.Validate, base *logrus.Logger, deviceChecker DeviceChecker) WorkoutDataService {
	return &WorkoutDataServiceImpl{
		workoutDataRepository: workoutDataRepository,
		validate: validate,
		log: helper.NewModuleLogger(base, "service", "workout_data"),
		deviceChecker: deviceChecker,
	}
}

func (s *WorkoutDataServiceImpl) Create(c fiber.Ctx, request WorkoutDataCreateRequest) (*WorkoutDataResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("request", request)
	
	log.Debug("received create workout data request")
	
	err := s.validate.Struct(request)
	if err != nil {
		log.WithError(err).Warn("failed to validate create workout data request")
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
	
	workoutData := &WorkoutData{
		Base: model.Base{ID: uuid.New()},
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

	workoutData, err = s.workoutDataRepository.Save(c, tx, workoutData)
	if err != nil {
		log.WithError(err).Error("failed to save workout data")
		return nil, fmt.Errorf("save workout data: %w", err)
	}

	log.WithField("workout_data_id", workoutData.ID).Debug("workout data created successfully")

	return toWorkoutDataResponse(workoutData), nil
}

func (s *WorkoutDataServiceImpl) GetByDeviceID(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("device_id", deviceID)
	
	log.Debug("received find workout data by device_id request")

	if deviceID == "" {
		log.Warn("validation failed: device_id required")
		return nil, &apperror.ValidationError{Message: "device_id required"}
	}
	
	authUser, err := context.GetAuthUser(c)
	if err != nil {
		log.WithError(err).Warn("user not authenticated")
		return nil, err
	}
	log = log.WithField("user_id", authUser.UserID)

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)
	
	owned, err := s.deviceChecker.IsOwnedByUser(c, tx, deviceID, authUser.UserID)
	if err != nil {
		log.WithError(err).Error("failed to check device ownership")
		return nil, err
	}
	if !owned {
		log.Warn("access denied: device belongs to another user")
		return nil, &apperror.ForbiddenError{Message: "This device belongs to another user"}
	}

	workoutDatas, err := s.workoutDataRepository.FindByDeviceID(c, tx, deviceID)
	if err != nil {
		log.WithError(err).Error("failed to find workout data by device_id")
		return nil, fmt.Errorf("find workout data by device_id %s: %w", deviceID, err)
	}

	log.WithField("result_count", len(workoutDatas)).Debug("workout data found successfully")

	return response.ToResponses(&workoutDatas, toWorkoutDataResponse), nil
}

func toWorkoutDataResponse(workoutData *WorkoutData) *WorkoutDataResponse {
	return &WorkoutDataResponse{
		WorkoutDataId:   workoutData.Base.ID.String(),
		WorkoutDataType: workoutData.WorkoutDataType,
		TotalSteps: workoutData.TotalSteps,
		TotalDistance: workoutData.TotalDistance,
		TotalCalories: workoutData.TotalCalories,
		HeartRateAverage: workoutData.HeartRateAverage,
		Pace: workoutData.Pace,
		Duration: workoutData.Duration,
		CreatedAt: workoutData.Base.CreatedAt.String(),
		DeviceID: workoutData.DeviceId,
	}
}

