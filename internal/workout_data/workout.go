package workout

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("workout not found")
	ErrDuplicate = errors.New("workout already uploaded")
)

type Type string
type Kind string

const (
	TypeRunning Type = "running"
	TypeWalking Type = "walking"
)

const (
	TypeDefault  Kind = "default"
	TypeWarmUp   Kind = "warm_up"
	TypeFastest  Kind = "fastest"
	TypeCoolDown Kind = "cool_down"
)

func (t Type) IsValid() bool { return t == TypeRunning || t == TypeWalking }

func (k Kind) IsValid() bool {
	return k == TypeDefault || k == TypeWarmUp || k == TypeFastest || k == TypeCoolDown
}

type Workout struct {
	ID                uuid.UUID
	ClientSessionID   string
	UserID            uuid.UUID
	DeviceName        string
	DeviceMacAddress  string
	Type              Type
	StartedAt         time.Time
	ActiveDurationSec int
	Steps             int
	DistanceM         float64
	CaloriesKcal      float64
	AvgPaceSecPerKm   int
	AvgHeartRateBpm   int
	MaxHeartRateBpm   int
	HeartRates        []HeartRatePoint
	Splits            []PaceSplit
	CreatedAt         time.Time
}

type HeartRatePoint struct {
	ActiveMs time.Time
	Bpm      int
}

type PaceSplit struct {
	PaceSecPerKm    int
	AvgHeartRateBpm int
	Kind            Kind
}
