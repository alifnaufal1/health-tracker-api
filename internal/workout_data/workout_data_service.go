package workout_data

import (
	"fmt"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/model"
	"health-tracker-api/pkg/response"
	"slices"
	"time"

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
	GetAll(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error)
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

	authUser, err := context.GetAuthUser(c)
	if err != nil {
		log.WithError(err).Warn("user not authenticated")
		return nil, err
	}
	log = log.WithField("user_id", authUser.UserID)
	
	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	owned, err := s.deviceChecker.IsOwnedByUser(c, tx, request.DeviceID, authUser.UserID)
	if err != nil {
		log.WithError(err).Error("failed to check device ownership")
		return nil, err
	}
	if !owned {
		log.Warn("access denied: device belongs to another user")
		return nil, &apperror.ForbiddenError{Message: "This device belongs to another user"}
	}
	
	totalSteps := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalSteps - request.WorkoutBatchData[0].TotalSteps
	totalDistance := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalDistance - request.WorkoutBatchData[0].TotalDistance
	totalCalories := request.WorkoutBatchData[len(request.WorkoutBatchData)-1].TotalCalories - request.WorkoutBatchData[0].TotalCalories
	heartRateAvg := 0.0
	var heartRateList []float64
	var paceList []float64
	var heartRateOvertime []HeartRateOverTime
	var dataOvertime []helper.DataOvertime
	for _, data := range request.WorkoutBatchData {
		heartRateAvg += data.HeartRate
		heartRateList = append(heartRateList, data.HeartRate)
		paceList = append(paceList, data.Pace)
		heartRateOvertime = append(heartRateOvertime, HeartRateOverTime{
			HeartRate: data.HeartRate,
			CreatedAt: data.CreatedAt,
		})
		dataOvertime = append(dataOvertime, helper.DataOvertime{
			HeartRate: data.HeartRate,
			TotalDistance: data.TotalDistance,
			CreatedAt: data.CreatedAt,
		})
	}
	heartRateAvg /= float64(len(request.WorkoutBatchData))
	heartRateMax := slices.Max(heartRateList)
	
	paceAvg, err, duration := helper.CountPace(request.WorkoutBatchData[0].CreatedAt, request.EndedAt, &totalDistance)
	if err != nil {
		log.WithError(err).Error("failed to count pace")
		return nil, err
	}
	paceMax := slices.Max(paceList)
	detailPerKm, err := helper.CountPacePerKm(&totalDistance, dataOvertime)
	if err != nil {
		log.WithError(err).Error("failed to count pace per km")
		return nil, err
	}
	fmt.Println("detailPerKm", detailPerKm)
	var finalDetailPerKm []DetailPerKm
	for _, data := range detailPerKm {
		finalDetailPerKm = append(finalDetailPerKm, DetailPerKm{
			PaceAvg: data.PaceAvg,
			HeartRateAvg: data.HeartRateAvg,
		})
	}

	endedAt, err := time.Parse(time.RFC3339, request.EndedAt)
	if err != nil {
		log.WithError(err).Error("failed to parse ended_at")
		return nil, err
	}

	workoutData := &WorkoutData{
		Base: model.Base{ID: uuid.New()},
		WorkoutDataType: request.WorkoutDataType,
		TotalSteps: totalSteps,
		TotalDistance: totalDistance,
		TotalCalories: totalCalories,
		HeartRateAvg: heartRateAvg,
		HeartRateMax: heartRateMax,
		HeartRateOverTime: heartRateOvertime,
		PaceAvg: paceAvg,
		PaceMax: paceMax,
		DetailPerKm: finalDetailPerKm,
		EndedAt: endedAt,
		Duration: *duration,
		DeviceId: request.DeviceID,
	}

	workoutDataResult, err := s.workoutDataRepository.Save(c, tx, workoutData)
	if err != nil {
		log.WithError(err).Error("failed to save workout data")
		return nil, fmt.Errorf("save workout data: %w", err)
	}

	log.WithField("workout_data_id", workoutDataResult.ID).Debug("workout data created successfully")

	return toWorkoutDataResponse(workoutData), nil
}

func (s *WorkoutDataServiceImpl) GetAll(c fiber.Ctx, deviceID string) (*[]WorkoutDataResponse, error) {
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
		HeartRateAvg: workoutData.HeartRateAvg,
		HeartRateMax: workoutData.HeartRateMax,
		HeartRateOverTime: workoutData.HeartRateOverTime,
		PaceAvg: workoutData.PaceAvg,
		PaceMax: workoutData.PaceMax,
		DetailPerKm: workoutData.DetailPerKm,
		CreatedAt: workoutData.Base.CreatedAt.String(),
		EndedAt: workoutData.EndedAt.String(),
		Duration: workoutData.Duration,
		DeviceID: workoutData.DeviceId,
	}
}

