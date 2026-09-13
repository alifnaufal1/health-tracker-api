package wokoutdata

type WorkoutDataCreateRequest struct {
	WorkoutDataId    string `json:"workout_data_id"`
	WorkoutDataType  string `json:"workout_data_type"`
	TotalSteps       int    `json:"total_steps"`
	TotalDistance    int    `json:"total_distance"`
	TotalCalories    int    `json:"total_calories"`
	HeartRateAverage int    `json:"heart_rate_average"`
	Pace             int    `json:"pace"`
	CreatedAt        string `json:"created_at"`
	TotalTime        string `json:"total_time"`
	DeviceID         string `json:"service_id"`
}