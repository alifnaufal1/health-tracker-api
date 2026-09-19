package workout_data

import (
	"errors"
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
)

type DeviceChecker interface {
	IsOwnedByUser(c fiber.Ctx, deviceID string, userID uuid.UUID) (bool, error)
}

type WorkoutDataService interface {
	Create(c fiber.Ctx, request WorkoutDataCreateRequest) (*WorkoutDataResponse, error)
	// Delete(c fiber.Ctx, workoutDataID string) error
	// FindByID(c fiber.Ctx, workoutDataID string) (*WorkoutDataResponse, error)
	GetByDeviceID(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error)
	// FindByAll(c fiber.Ctx) (*[]WorkoutDataResponse, error)
}

type WorkoutDataServiceImpl struct {
	workoutDataRepository WorkoutDataRepository
	validate *validator.Validate
	log *logrus.Logger
	deviceChecker DeviceChecker 
}

func NewWorkoutDataService(workoutDataRepository WorkoutDataRepository, validate *validator.Validate, log *logrus.Logger, deviceChecker DeviceChecker) WorkoutDataService {
	return &WorkoutDataServiceImpl{
		workoutDataRepository: workoutDataRepository,
		validate: validate,
		log: log,
		deviceChecker: deviceChecker,
	}
}

func (s *WorkoutDataServiceImpl) Create(c fiber.Ctx, request WorkoutDataCreateRequest) (*WorkoutDataResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)
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
		log.Warn("Create workout data process failed at layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("workout_data_id", workoutData.ID).Info("workout data created successfully")

	return toWorkoutDataResponse(workoutData), nil
}

func (s *WorkoutDataServiceImpl) GetByDeviceID(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("device_id", deviceID)
	
	log.Info("Received find workout data by device_id request")

	if deviceID == "" {
		log.Warn("Validation failed for find workout data by device_id request")
		return nil, &apperror.ValidationError{Message: "device_id required"}
	}
	
	authUser, err := context.GetAuthUser(c)
	if err != nil {
		log.WithField("authUser", authUser).Warn("User not authenticated")
		return nil, err
	}

	owned, err := s.deviceChecker.IsOwnedByUser(c, deviceID, authUser.UserID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, &apperror.ForbiddenError{Message: "This workout data is belongs to another device"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	workout_datas, err := s.workoutDataRepository.FindByDeviceID(c, tx, deviceID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.WithField("workout_datas", workout_datas).Info("workout data found successfully")

	return response.ToResponses(&workout_datas, toWorkoutDataResponse), nil
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

