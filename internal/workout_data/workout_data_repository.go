package workout_data

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type WorkoutDataRepository interface {
	Save(c fiber.Ctx, tx *gorm.DB, workoutData *WorkoutData) (*WorkoutData, error)
	Delete(c fiber.Ctx, tx *gorm.DB, workoutDataID string) error
	FindByID(c fiber.Ctx, tx *gorm.DB, workoutDataID string) (*WorkoutData, error)
	FindByDeviceID(c fiber.Ctx, tx *gorm.DB, deviceID string) ([]WorkoutData, error)
	FindAll(c fiber.Ctx, tx *gorm.DB) (*[]WorkoutData, error)
}

type WorkoutDataRepositoryImpl struct{}

func NewWorkoutDataRepository() WorkoutDataRepository {
	return &WorkoutDataRepositoryImpl{}
}

func (r *WorkoutDataRepositoryImpl) Save(c fiber.Ctx, tx *gorm.DB, workoutData *WorkoutData) (*WorkoutData, error) {
	result := gorm.WithResult()
	err := gorm.G[WorkoutData](tx, result).Create(c, workoutData)
	if err != nil {
		return nil, err
	}
	return workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) Delete(c fiber.Ctx, tx *gorm.DB, workoutDataID string) error {
	result := gorm.WithResult()
	_, err := gorm.G[WorkoutData](tx, result).Where("id = ?", workoutDataID).Delete(c)
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkoutDataRepositoryImpl) FindByID(c fiber.Ctx, tx *gorm.DB, workoutDataID string) (*WorkoutData, error) {
	result := gorm.WithResult()
	workoutData, err := gorm.G[WorkoutData](tx, result).Where("id = ?", workoutDataID).First(c)
	if err != nil {
		return nil, err
	}
	return &workoutData, nil
}

func (r *WorkoutDataRepositoryImpl) FindByDeviceID(c fiber.Ctx, tx *gorm.DB, deviceID string) ([]WorkoutData, error) {
	result := gorm.WithResult()
	workoutDatas := make([]WorkoutData, 0)
	workoutDatas, err := gorm.G[WorkoutData](tx, result).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").Find(c)
	if err != nil {
		return nil, err
	}
	return workoutDatas, nil
}

func (r *WorkoutDataRepositoryImpl) FindAll(c fiber.Ctx, tx *gorm.DB) (*[]WorkoutData, error) {
	result := gorm.WithResult()
	workoutDatas, err := gorm.G[WorkoutData](tx, result).Find(c)
	if err != nil {
		return nil, err
	}
	return &workoutDatas, nil
}
