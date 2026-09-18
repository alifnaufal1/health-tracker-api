package workout_data

type WorkoutDataCreateRequest struct {
	WorkoutDataType  string             `validate:"required" json:"workout_data_type"`
	DeviceID         string             `validate:"required" json:"device_id"`
	Duration         int64              `validate:"required" json:"duration"`
	Pace             int                `validate:"required" json:"pace"`
	WorkoutBatchData []WorkoutBatchData `validate:"required" json:"workout_batch_data"`
}

type WorkoutDataResponse struct {
	WorkoutDataId    string `json:"workout_data_id"`
	WorkoutDataType  string `json:"workout_data_type"`
	TotalSteps       int    `json:"total_steps"`
	TotalDistance    int    `json:"total_distance"`
	TotalCalories    int    `json:"total_calories"`
	HeartRateAverage int    `json:"heart_rate_average"`
	Pace             int    `json:"pace"`
	CreatedAt        string `json:"created_at"`
	Duration         int64  `json:"duration"`
	DeviceID         string `json:"service_id"`
}