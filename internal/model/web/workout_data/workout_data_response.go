package web

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