package repository

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WorkoutDataRepository interface {
	Save(ctx fiber.Ctx, tx *gorm.DB, workoutData *domain.WorkoutData) (*domain.WorkoutData, error)
	Delete(ctx fiber.Ctx, tx *gorm.DB, workoutDataId string) error
	FindById(ctx fiber.Ctx, tx *gorm.DB, workoutDataId string) (*domain.WorkoutData, error)
	FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.WorkoutData, error)
}

type WorkoutDataRepositoryImpl struct {
	log *logrus.Logger
}

func NewWorkoutDataRepository(log *logrus.Logger) WorkoutDataRepository {
	return &WorkoutDataRepositoryImpl{log: log}
}

func (r *WorkoutDataRepositoryImpl) Save(ctx fiber.Ctx, tx *gorm.DB, workoutData *domain.WorkoutData) (*domain.WorkoutData, error) {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("Workout data", workoutData)

	log.Debug("Inserting new workout data into database")

	result := gorm.WithResult()
	err := gorm.G[domain.WorkoutData](tx, result).Create(ctx, workoutData)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database insert failed")
		return nil, err
	}

	log.WithField("workout_data_id", workoutData.ID).Debug("Workout data inserted successfully")

	return workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) Delete(ctx fiber.Ctx, tx *gorm.DB, workoutDataID string) error {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("workout_data_id", workoutDataID)
	
	log.Debug("Deleting workout data from database")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.WorkoutData](tx, result).Where("id = ?", workoutDataID).Delete(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Workout data delete failed")
		return err
	}
	
	log.WithField("workout_data_id", workoutDataID).Debug("Workout data deleted successfully")
	
	return nil
}

func (r *WorkoutDataRepositoryImpl) FindById(ctx fiber.Ctx, tx *gorm.DB, workoutDataID string) (*domain.WorkoutData, error) {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("workout_data_id", workoutDataID)

	log.Debug("Finding workout data by id from database")
	
	result := gorm.WithResult()
	workoutData, err := gorm.G[domain.WorkoutData](tx, result).Where("id = ?", workoutDataID).First(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Workout data find by id failed")
		return nil, err
	}
	
	log.WithField("workout_data_id", workoutDataID).Debug("Workout data found by id successfully")
	
	return &workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.WorkoutData, error) {
	log := helper.LoggerWithRequestID(ctx, r.log)

	log.Debug("Finding all workout data from database")
	
	result := gorm.WithResult()
	workoutDatas, err := gorm.G[domain.WorkoutData](tx, result).Find(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find all failed")
		return nil, err
	}
	
	log.Debug("All workout data found successfully")
	
	return &workoutDatas, nil
}