package helper

import (
	"fmt"
	"time"
)

type DataOvertime struct {
	HeartRate     float64
	TotalDistance float64
	CreatedAt     string
}

type DetailPerKm struct {
	PaceAvg      float64    
	HeartRateAvg float64
}

func CountPace(isoStartTime, isoEndTime string, totalDistance *float64) (float64, error, *int64) {
	totalDistanceInKm := *totalDistance / 1000
	startTime, err := time.Parse(time.RFC3339, isoStartTime)
	if err != nil {	
		return 0, err, nil
	}	
	fmt.Println("startTime", startTime)
	
	endTime, err := time.Parse(time.RFC3339, isoEndTime)
	if err != nil {
		return 0, err, nil
	}
	fmt.Println("endTime", endTime)
	
	duration := endTime.Unix() - startTime.Unix()
	fmt.Println("duration", duration)
	durationInMinute := float64(int(duration) / 60)
	fmt.Println("durationInMinute", durationInMinute)
	pace := durationInMinute / totalDistanceInKm
	fmt.Println("pace", pace)

	return pace, nil, &duration
}

func CountPacePerKm(totalDistance *float64, dataOvertime []DataOvertime) ([]DetailPerKm, error) {
	var groupData []DetailPerKm = []DetailPerKm{}
	totalKm := int((*totalDistance + 1000 - 1) / 1000)
	firstIndex := 0

	for i := 1; i <= totalKm; i++ {
		distance := i * 1000
		for j, data := range dataOvertime {
			totalHeartRate := 0.0
			currentDistance := data.TotalDistance - dataOvertime[firstIndex].TotalDistance
			if int(currentDistance) < distance{
				totalHeartRate += data.HeartRate
			} else {
				firstIndex = j + 1
				groupData[i-1].HeartRateAvg= totalHeartRate / float64(firstIndex) //i'm not sure, because heart rate not unique (there are more than one time that value is same)
				paceAvg, err, _ := CountPace(dataOvertime[0].CreatedAt, data.CreatedAt, &currentDistance)
				if err != nil {
					return nil, err
				}
				groupData[i-1].PaceAvg = paceAvg
				break
			}
		}
	}

	return groupData, nil
}