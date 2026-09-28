package helper

import (
	"math"
	"time"
)

type DataOvertime struct {
	HeartRate     float64
	TotalDistance float64
	CreatedAt     string
}
type DetailPerKm struct {
	AvgPace      int
	AvgHeartRate float64
}

func CountPace(isoStartTime, isoEndTime *string, totalDistance *float64) (*int, *int64, error) {
	totalDistanceInKm := *totalDistance / 1000

	startTime, err := time.Parse(time.RFC3339, *isoStartTime)
	if err != nil {
		return nil, nil, err
	}

	endTime, err := time.Parse(time.RFC3339, *isoEndTime)
	if err != nil {
		return nil, nil, err
	}

	durationInSeconds := endTime.Unix() - startTime.Unix()
	pace := float64(durationInSeconds) / totalDistanceInKm
	pace = math.Round(pace*100) / 100
	paceSecPerKm := int(pace)
	return &paceSecPerKm, &durationInSeconds, nil
}

func CountPacePerKm(totalDistance *float64, dataOvertime []DataOvertime) ([]DetailPerKm, error) {
	var groupData []DetailPerKm
	totalKm := int((*totalDistance + 1000 - 1) / 1000)
	firstIndex := 0

	for i := 1; i <= totalKm; i++ {
		distance := i * 1000
		totalHeartRate := 0.0
		for j, data := range dataOvertime {
			currentDistance := data.TotalDistance - dataOvertime[firstIndex].TotalDistance
			if int(currentDistance) >= distance || j+1 == len(dataOvertime) {
				totalHeartRate += data.HeartRate
				firstIndex = j + 1
				avgPace, _, err := CountPace(&dataOvertime[0].CreatedAt, &data.CreatedAt, &currentDistance)
				if err != nil {
					return nil, err
				}
				groupData = append(groupData, DetailPerKm{
					AvgHeartRate: totalHeartRate / float64(firstIndex), //i'm not sure, because heart rate not unique (there are more than one time that value is same)
					AvgPace:      *avgPace,
				})
				break
			} else {
				totalHeartRate += data.HeartRate
			}
		}
	}
	return groupData, nil
}

// func CountPace(isoStartTime, isoEndTime *string, totalDistance *float64) (string, error, *int64) {
// 	totalDistanceInKm := *totalDistance / 1000

// 	startTime, err := time.Parse(time.RFC3339, *isoStartTime)
// 	if err != nil {
// 		return "", err, nil
// 	}

// 	endTime, err := time.Parse(time.RFC3339, *isoEndTime)
// 	if err != nil {
// 		return "", err, nil
// 	}

// 	durationInSeconds := endTime.Unix() - startTime.Unix()
// 	durationInMinute := float64(durationInSeconds) / 60
// 	pace := durationInMinute / totalDistanceInKm
// 	pace = math.Round(pace*100) / 100

// 	intPart, fracPart := math.Modf(pace)
// 	var finalPace string

// 	if fracPart == float64(0) {
// 		finalPace = fmt.Sprintf("%02d:%02d", int(intPart), int(fracPart))
// 	} else {
// 		second := fracPart * 60
// 		deepIntPart, _ := math.Modf(second)
// 		finalPace = fmt.Sprintf("%02d:%02d", int(intPart), int(deepIntPart))
// 	}

// 	return finalPace, nil, &durationInSeconds
// }
