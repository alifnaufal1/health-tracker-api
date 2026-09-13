package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutData struct {
	Base
	WorkoutDataType  string `gorm:"not null"`
	TotalSteps       int    `gorm:"not null"`
	TotalDistance    int    `gorm:"not null"`
	TotalCalory      int    `gorm:"not null"`
	HeartRateAverage int    `gorm:"not null"`
	Pace             int    `gorm:"not null"`
	WorkoutBatchData []WorkoutBatchData
	DeviceId         *uuid.UUID
	Device 			 Device
}

type WorkoutBatchData struct {
	TotalSteps       int    
	TotalDistance    int    
	TotalCalories      int    
	HeartRateAverage int    
	Pace             int    
	CreatedAt		 time.Time
}