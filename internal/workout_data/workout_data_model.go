package workout_data

import (
	"health-tracker-api/internal/device"
	"health-tracker-api/pkg/model"
	"time"

	"gorm.io/datatypes"
)

type WorkoutData struct {
	model.Base
	WorkoutDataType   string                                 `gorm:"not null"`
	TotalSteps        float64                                `gorm:"not null"`
	TotalDistance     float64                                `gorm:"not null"`
	TotalCalories     float64                                `gorm:"not null"`
	HeartRateAvg      float64                                `gorm:"not null"`
	HeartRateMax      float64                                `gorm:"not null"`
	HeartRateOverTime datatypes.JSONSlice[HeartRateOverTime] `gorm:"not null"`
	PaceAvg           string                                 `gorm:"not null"`
	PaceMax           string                                 `gorm:"not null"`
	DetailPerKm       datatypes.JSONSlice[DetailPerKm]       `gorm:"not null"`
	EndedAt           time.Time                              `gorm:"not null"`
	Duration          int64                                  `gorm:"not null"`
	DeviceId          string
	Device            device.Device
}

type HeartRateOverTime struct {
	HeartRate float64 `json:"heart_rate"`
	CreatedAt string  `json:"created_at"`
}

type DetailPerKm struct {
	PaceAvg      string  `json:"pace_avg"`
	HeartRateAvg float64 `json:"heart_rate_avg"`
}
