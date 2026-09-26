package helper

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type DataOvertime struct {
	HeartRate     float64
	TotalDistance float64
	CreatedAt     string
}
type DetailPerKm struct {
	PaceAvg      string    
	HeartRateAvg float64
}

func CountPace(isoStartTime, isoEndTime *string, totalDistance *float64) (string, error, *int64) {
	totalDistanceInKm := *totalDistance / 1000

	startTime, err := time.Parse(time.RFC3339, *isoStartTime)
	if err != nil {	
		return "", err, nil
	}	
	
	endTime, err := time.Parse(time.RFC3339, *isoEndTime)
	if err != nil {
		return "", err, nil
	}
	
	durationInSeconds := endTime.Unix() - startTime.Unix()
	durationInMinute := float64(durationInSeconds) / 60	
	pace := durationInMinute / totalDistanceInKm
	pace = math.Round(pace*100) / 100

	intPart, fracPart := math.Modf(pace)
	var finalPace string

	if fracPart == float64(0) {
		finalPace = fmt.Sprintf("%02d:%02d", int(intPart), int(fracPart))
	} else {
		second := fracPart * 60
		deepIntPart, _ := math.Modf(second)
		finalPace = fmt.Sprintf("%02d:%02d", int(intPart), int(deepIntPart))
	}

	return finalPace, nil, &durationInSeconds
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
				paceAvg, err, _ := CountPace(&dataOvertime[0].CreatedAt, &data.CreatedAt, &currentDistance)
				if err != nil {
					return nil, err
				}
				groupData = append(groupData, DetailPerKm{
					HeartRateAvg: totalHeartRate / float64(firstIndex), //i'm not sure, because heart rate not unique (there are more than one time that value is same)
					PaceAvg: paceAvg,
				})
				break
			} else {
				totalHeartRate += data.HeartRate
			}
		}
	}
	return groupData, nil
}

func FindPaceMax(paceList []string) (*string, error) {
	var rawPaceList []struct{
		rawPace float64
		index int
	}

	for i, pace := range paceList {
		time := strings.Split(pace, ":")
		if len(time) != 2 {
			continue
		}
		currentMinute, err := strconv.Atoi(time[0])
		if err != nil {
			return nil, errors.New("the pace format is incorrect")
		}
		currentSecond, err := strconv.Atoi(time[1])
		if err != nil {
			return nil, errors.New("the pace format is incorrect")
		}
		rawPace := float64(currentMinute) + (float64(currentSecond)/60)
		rawPaceList = append(rawPaceList, struct{rawPace float64; index int}{
			rawPace: rawPace,
			index: i,
		})
	}

	maxPace := slices.MinFunc(rawPaceList, func(a, b struct{rawPace float64; index int}) int {
		return cmp.Compare(a.rawPace, b.rawPace)
	})
	return &paceList[maxPace.index], nil
}