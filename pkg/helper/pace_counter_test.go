package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountPace(t *testing.T) {
	tests := []struct {
		testId                    string
		isoStart                  string
		isoEnd                    string
		totalDistanceInM          float64
		expectedPace              string
		expectedDurationInSeconds int64
	}{
		{
			testId:                    "Test 5 km in 30 minutes",
			isoStart:                  "2026-09-14T08:00:00Z",
			isoEnd:                    "2026-09-14T08:30:00Z",
			totalDistanceInM:          5000,
			expectedPace:              "06:00",
			expectedDurationInSeconds: 1800,
		},
		{
			testId:                    "Test 500 m in 5 minutes",
			isoStart:                  "2026-09-14T09:00:00Z",
			isoEnd:                    "2026-09-14T09:05:00Z",
			totalDistanceInM:          500,
			expectedPace:              "10:00",
			expectedDurationInSeconds: 300,
		},
		{
			testId:                    "Test 900 m in 3 minutes 30 seconds",
			isoStart:                  "2026-09-14T10:00:00Z",
			isoEnd:                    "2026-09-14T10:03:30Z",
			totalDistanceInM:          900,
			expectedPace:              "03:53",
			expectedDurationInSeconds: 210,
		},
		{
			testId:                    "Test 100 m in 20 seconds",
			isoStart:                  "2026-09-14T11:00:00Z",
			isoEnd:                    "2026-09-14T11:00:20Z",
			totalDistanceInM:          100,
			expectedPace:              "03:20",
			expectedDurationInSeconds: 20,
		},
		{
			testId:                    "Test 50 m in 1 minute",
			isoStart:                  "2026-09-14T12:00:00Z",
			isoEnd:                    "2026-09-14T12:01:00Z",
			totalDistanceInM:          50,
			expectedPace:              "20:00",
			expectedDurationInSeconds: 60,
		},
		{
			testId:                    "Test 10 m in 10 seconds",
			isoStart:                  "2026-09-14T13:00:00Z",
			isoEnd:                    "2026-09-14T13:00:10Z",
			totalDistanceInM:          10,
			expectedPace:              "16:40",
			expectedDurationInSeconds: 10,
		},
		{
			testId:                    "Test 99 m in 5 seconds",
			isoStart:                  "2026-09-14T14:00:00Z",
			isoEnd:                    "2026-09-14T14:00:05Z",
			totalDistanceInM:          99,
			expectedPace:              "00:50",
			expectedDurationInSeconds: 5,
		},
		{
			testId:                    "Test 50 m in 30 seconds",
			isoStart:                  "2026-09-14T15:00:00Z",
			isoEnd:                    "2026-09-14T15:00:30Z",
			totalDistanceInM:          50,
			expectedPace:              "10:00",
			expectedDurationInSeconds: 30,
		},
		{
			testId:                    "Test 1 m in 1 second",
			isoStart:                  "2026-09-14T16:00:00Z",
			isoEnd:                    "2026-09-14T16:00:01Z",
			totalDistanceInM:          1,
			expectedPace:              "16:40",
			expectedDurationInSeconds: 1,
		},
		{
			testId:                    "Test 0.5 m in 10 seconds",
			isoStart:                  "2026-09-14T17:00:00Z",
			isoEnd:                    "2026-09-14T17:00:10Z",
			totalDistanceInM:          0.5,
			expectedPace:              "333:20",
			expectedDurationInSeconds: 10,
		},
		{
			testId:                    "Test 0.99 m in 3 seconds",
			isoStart:                  "2026-09-14T18:00:00Z",
			isoEnd:                    "2026-09-14T18:00:03Z",
			totalDistanceInM:          0.99,
			expectedPace:              "50:30",
			expectedDurationInSeconds: 3,
		},
		{
			testId:                    "Test 0.1 m in 1 second",
			isoStart:                  "2026-09-14T19:00:00Z",
			isoEnd:                    "2026-09-14T19:00:01Z",
			totalDistanceInM:          0.1,
			expectedPace:              "166:40",
			expectedDurationInSeconds: 1,
		},
		{
			testId:                    "Test 250 m in 12 minutes 45 seconds",
			isoStart:                  "2026-09-14T20:00:00Z",
			isoEnd:                    "2026-09-14T20:12:45Z",
			totalDistanceInM:          250,
			expectedPace:              "51:00",
			expectedDurationInSeconds: 765,
		},
		{
			testId:                    "Test 750 m in 8 minutes 13 seconds",
			isoStart:                  "2026-09-15T08:00:00Z",
			isoEnd:                    "2026-09-15T08:08:13Z",
			totalDistanceInM:          750,
			expectedPace:              "10:57",
			expectedDurationInSeconds: 493,
		},
		{
			testId:                    "Test 1.5 km in 11 minutes",
			isoStart:                  "2026-09-15T09:00:00Z",
			isoEnd:                    "2026-09-15T09:11:00Z",
			totalDistanceInM:          1500,
			expectedPace:              "07:20",
			expectedDurationInSeconds: 660,
		},
		{
			testId:                    "Test 2.75 km in 17 minutes 29 seconds",
			isoStart:                  "2026-09-15T10:00:00Z",
			isoEnd:                    "2026-09-15T10:17:29Z",
			totalDistanceInM:          2750,
			expectedPace:              "06:21",
			expectedDurationInSeconds: 1049,
		},
		{
			testId:                    "Test 3.2 km in 25 minutes 37 seconds",
			isoStart:                  "2026-09-15T11:00:00Z",
			isoEnd:                    "2026-09-15T11:25:37Z",
			totalDistanceInM:          3200,
			expectedPace:              "08:01",
			expectedDurationInSeconds: 1537,
		},
		{
			testId:                    "Test 7.25 km in 41 minutes 11 seconds",
			isoStart:                  "2026-09-15T12:00:00Z",
			isoEnd:                  "2026-09-15T12:41:11Z",
			totalDistanceInM:          7250,
			expectedPace:              "05:41",
			expectedDurationInSeconds: 2471,
		},
		{
			testId:                    "Test 12.5 km in 1 hour 7 minutes",
			isoStart:                  "2026-09-15T13:00:00Z",
			isoEnd:                  "2026-09-15T14:07:00Z",
			totalDistanceInM:          12500,
			expectedPace:              "05:22",
			expectedDurationInSeconds: 4020,
		},
	}

	for _, test := range tests {
		t.Run(test.testId, func(t *testing.T) {
			pace, _, duration := CountPace(&test.isoStart, &test.isoEnd, &test.totalDistanceInM)
			assert.Equal(t, test.expectedPace, pace)
			assert.Equal(t, test.expectedDurationInSeconds, *duration)
		})
	}
}
