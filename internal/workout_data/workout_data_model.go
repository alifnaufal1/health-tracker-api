package workout_data

import (
	"health-tracker-api/internal/device"
	"health-tracker-api/pkg/model"
	"time"

	"gorm.io/datatypes"
)

type WorkoutData struct {
	model.Base
	WorkoutDataType  string                                `gorm:"not null"`
	TotalSteps       int                                   `gorm:"not null"`
	TotalDistance    int                                   `gorm:"not null"`
	TotalCalories    int                                   `gorm:"not null"`
	HeartRateAverage int                                   `gorm:"not null"`
	Pace             int                                   `gorm:"not null"`
	Duration 		 int64 								   `gorm:"not null"`
	WorkoutBatchData datatypes.JSONSlice[WorkoutBatchData] `gorm:"type:jsonb;not null"`
	DeviceId         string
	Device           device.Device
}

type WorkoutBatchData struct {
	TotalSteps       int `validate:"required" json:"total_steps"`
	TotalDistance    int `validate:"required" json:"total_distance"`
	TotalCalories    int `validate:"required" json:"total_calories"`
	HeartRate        int `validate:"required" json:"heart_rate"`
	CreatedAt        time.Time `validate:"required" json:"created_at"`
}