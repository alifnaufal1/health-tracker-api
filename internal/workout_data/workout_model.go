package workout

import (
	"health-tracker-api/internal/device"
	"health-tracker-api/pkg/model"
	"time"

	"gorm.io/datatypes"
)

type WorkoutData struct {
	model.Base
	WorkoutDataType string                              `gorm:"not null"`
	TotalSteps      float64                             `gorm:"not null"`
	TotalDistance   float64                             `gorm:"not null"`
	TotalCalories   float64                             `gorm:"not null"`
	AvgHeartRate    float64                             `gorm:"not null"`
	MaxHeartRate    float64                             `gorm:"not null"`
	AvgPace         int                                 `gorm:"not null"`
	BestPace        int                                 `gorm:"not null"`
	Duration        int64                               `gorm:"not null"`
	StartedAt       time.Time                           `gorm:"not null"`
	EndedAt         time.Time                           `gorm:"not null"`
	HeartRateSeries datatypes.JSONSlice[HeartRatePoint] `gorm:"not null"`
	Splits          datatypes.JSONSlice[PaceSplit]      `gorm:"not null"`
	DeviceId        string
	Device          device.Device
}
