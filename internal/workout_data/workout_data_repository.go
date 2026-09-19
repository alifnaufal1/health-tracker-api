package workout_data

import (
	"health-tracker-api/pkg/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WorkoutDataRepository interface {
	Save(c fiber.Ctx, tx *gorm.DB, workoutData *WorkoutData) (*WorkoutData, error)
	Delete(c fiber.Ctx, tx *gorm.DB, workoutDataID string) error
	FindByID(c fiber.Ctx, tx *gorm.DB, workoutDataID string) (*WorkoutData, error)
	FindByDeviceID(c fiber.Ctx, tx *gorm.DB, deviceID string) ([]WorkoutData, error)
	FindAll(c fiber.Ctx, tx *gorm.DB) (*[]WorkoutData, error)
}

type WorkoutDataRepositoryImpl struct {
	log *logrus.Logger
}

func NewWorkoutDataRepository(log *logrus.Logger) WorkoutDataRepository {
	return &WorkoutDataRepositoryImpl{log: log}
}

func (r *WorkoutDataRepositoryImpl) Save(c fiber.Ctx, tx *gorm.DB, workoutData *WorkoutData) (*WorkoutData, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("Workout data", workoutData)

	log.Debug("Inserting new workout data into database")

	result := gorm.WithResult()
	err := gorm.G[WorkoutData](tx, result).Create(c, workoutData)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database insert failed")
		return nil, err
	}

	log.WithField("workout_data_id", workoutData.ID).Debug("Workout data inserted successfully")

	return workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) Delete(c fiber.Ctx, tx *gorm.DB, workoutDataID string) error {
	log := helper.LoggerWithRequestID(c, r.log).WithField("workout_data_id", workoutDataID)
	
	log.Debug("Deleting workout data from database")
	
	result := gorm.WithResult()
	_, err := gorm.G[WorkoutData](tx, result).Where("id = ?", workoutDataID).Delete(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Workout data delete failed")
		return err
	}
	
	log.WithField("workout_data_id", workoutDataID).Debug("Workout data deleted successfully")
	
	return nil
}

func (r *WorkoutDataRepositoryImpl) FindByID(c fiber.Ctx, tx *gorm.DB, workoutDataID string) (*WorkoutData, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("workout_data_id", workoutDataID)

	log.Debug("Finding workout data by id from database")
	
	result := gorm.WithResult()
	workoutData, err := gorm.G[WorkoutData](tx, result).Where("id = ?", workoutDataID).First(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Workout data find by id failed")
		return nil, err
	}
	
	log.WithField("workout_data_id", workoutDataID).Debug("Workout data found by id successfully")
	
	return &workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) FindByDeviceID(c fiber.Ctx, tx *gorm.DB, deviceID string) ([]WorkoutData, error) {
	log := helper.LoggerWithRequestID(c, r.log).WithField("device_id", deviceID)

	log.Debug("Finding workout data by device_id from database")
	
	result := gorm.WithResult()
	workoutDatas := make([]WorkoutData, 0)
	workoutDatas, err := gorm.G[WorkoutData](tx, result).Where("device_id = ?", deviceID).Find(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Workout data find by device_id failed")
		return nil, err
	}
	
	log.WithField("device_id", deviceID).Debug("Workout data found by device_id successfully")
	
	return workoutDatas, nil
}

func (r *WorkoutDataRepositoryImpl) FindAll(c fiber.Ctx, tx *gorm.DB) (*[]WorkoutData, error) {
	log := helper.LoggerWithRequestID(c, r.log)

	log.Debug("Finding all workout data from database")
	
	result := gorm.WithResult()
	workoutDatas, err := gorm.G[WorkoutData](tx, result).Find(c)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find all failed")
		return nil, err
	}
	
	log.Debug("All workout data found successfully")
	
	return &workoutDatas, nil
}