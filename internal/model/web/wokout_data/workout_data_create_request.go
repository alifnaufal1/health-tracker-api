package wokoutdata

import "time"

type WorkoutDataCreateResponse struct {
	WorkoutDataType  string `validate:"required" json:"workout_data_type"`
	DeviceID         string `validate:"required" json:"service_id"`
	TotalTime        string `validate:"required" json:"total_time"`
	Pace             int    `validate:"required" json:"pace"`
	WorkoutBatchData []WorkoutBatchDataCreateResponse    `validate:"required" json:"workout_batch_data"`
}

type WorkoutBatchDataCreateResponse struct {
	TotalSteps       int `validate:"required" json:"total_steps"`
	TotalDistance    int `validate:"required" json:"total_distance"`
	TotalCalories    int `validate:"required" json:"total_calories"`
	HeartRate        int `validate:"required" json:"heart_rate"`
	CreatedAt        time.Time `validate:"required" json:"created_at"`
}