package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkoutData struct {
	Base
	WorkoutDataType  string                                `gorm:"not null"`
	TotalSteps       int                                   `gorm:"not null"`
	TotalDistance    int                                   `gorm:"not null"`
	TotalCalory      int                                   `gorm:"not null"`
	HeartRateAverage int                                   `gorm:"not null"`
	Pace             int                                   `gorm:"not null"`
	WorkoutBatchData datatypes.JSONSlice[WorkoutBatchData] `gorm:"type:jsonb;not null"`
	DeviceId         *uuid.UUID
	Device           Device
}

type WorkoutBatchData struct {
	TotalSteps       int       `json:"total_steps"`
	TotalDistance    int       `json:"total_distance"`
	TotalCalories    int       `json:"total_calories"`
	HeartRateAverage int       `json:"heart_rate_average"`
	Pace             int       `json:"pace"`
	CreatedAt        time.Time `json:"created_at"`
}