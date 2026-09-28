package workout_data

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

type HeartRatePoint struct {
	HeartRate float64 `json:"heart_rate"`
	Timestamp string  `json:"timestemp"`
}

type PaceSplit struct {
	Pace         int     `json:"pace"`
	AvgHeartRate float64 `json:"avg_heart_rate"`
	Type         float64 `json:"type"`
}
