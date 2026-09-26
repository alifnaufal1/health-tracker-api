package workout_data

type WorkoutDataCreateRequest struct {
	WorkoutDataType  string             `validate:"required" json:"workout_data_type"`
	DeviceID         string             `validate:"required" json:"device_id"`
	EndedAt          string             `validate:"required" json:"ended_at"`
	WorkoutBatchData []WorkoutBatchData `validate:"required" json:"workout_batch_data"`
}

type WorkoutBatchData struct {
	TotalSteps    float64 `validate:"required" json:"total_steps"`
	TotalDistance float64 `validate:"required" json:"total_distance"`
	TotalCalories float64 `validate:"required" json:"total_calories"`
	HeartRate     float64 `validate:"required" json:"heart_rate"`
	Pace          string  `validate:"required" json:"pace"`
	CreatedAt     string  `validate:"required" json:"created_at"`
}

type WorkoutDataResponse struct {
	WorkoutDataId     string              `json:"workout_data_id"`
	WorkoutDataType   string              `json:"workout_data_type"`
	TotalSteps        float64             `json:"total_steps"`
	TotalDistance     float64             `json:"total_distance"`
	TotalCalories     float64             `json:"total_calories"`
	HeartRateAvg      float64             `json:"heart_rate_avg"`
	HeartRateMax      float64             `json:"heart_rate_max"`
	HeartRateOverTime []HeartRateOverTime `json:"heart_rate_overtime"`
	PaceAvg           string              `json:"pace_avg"`
	PaceMax           string              `json:"pace_max"`
	DetailPerKm       []DetailPerKm       `json:"detail_per_km"`
	CreatedAt         string              `json:"created_at"`
	EndedAt           string              `json:"ended_at"`
	Duration          int64               `json:"duration"`
	DeviceID          string              `json:"device_id"`
}
