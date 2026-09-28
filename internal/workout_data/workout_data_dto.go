package workout_data

type WorkoutDataCreateRequest struct {
	WorkoutDataType  string             `validate:"required" json:"workout_data_type"`
	DeviceID         string             `validate:"required" json:"device_id"`
	WorkoutBatchData []WorkoutBatchData `validate:"required" json:"workout_batch_data"`
}

type WorkoutBatchData struct {
	Steps     float64 `validate:"required" json:"steps"`
	Distance  float64 `validate:"required" json:"distance"`
	Calories  float64 `validate:"required" json:"calories"`
	HeartRate float64 `validate:"required" json:"heart_rate"`
	Pace      int     `validate:"required" json:"pace"`
	Timestemp string  `validate:"required" json:"timestemp"`
}

type WorkoutDataResponse struct {
	WorkoutDataId   string           `json:"workout_data_id"`
	WorkoutDataType string           `json:"workout_data_type"`
	TotalSteps      float64          `json:"total_steps"`
	TotalDistance   float64          `json:"total_distance"`
	TotalCalories   float64          `json:"total_calories"`
	AvgHeartRate    float64          `json:"avg_heart_rate"`
	MaxHeartRate    float64          `json:"max_heart_rate"`
	AvgPace         int              `json:"avg_pace"`
	BestPace        int              `json:"best_pace"`
	Duration        int64            `json:"duration"`
	HeartRateSeries []HeartRatePoint `json:"heart_rate_series"`
	Splits          []PaceSplit      `json:"splits"`
	StartedAt       string           `json:"started_at"`
	EndedAt         string           `json:"ended_at"`
	DeviceID        string           `json:"device_id"`
}
