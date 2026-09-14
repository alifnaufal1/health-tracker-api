package web

import "health-tracker-api/internal/model/domain"

type WorkoutDataCreateRequest struct {
	WorkoutDataType  string                    `validate:"required" json:"workout_data_type"`
	DeviceID         string             	   `validate:"required" json:"device_id"`
	Duration         int64                	   `validate:"required" json:"duration"`
	Pace             int             		   `validate:"required" json:"pace"`
	WorkoutBatchData []domain.WorkoutBatchData `validate:"required" json:"workout_batch_data"`
}