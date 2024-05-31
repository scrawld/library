package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test for GetYearWeekByStartOfWeek
func TestGetYearWeekByStartOfWeek(t *testing.T) {
	tests := []struct {
		date         time.Time
		startOfWeek  time.Weekday
		expectedYear int
		expectedWeek int
	}{
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Monday, 2024, 22},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Monday, 2024, 1},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Sunday, 2023, 52},
	}

	for _, test := range tests {
		year, week := GetYearWeekByStartOfWeek(test.date, test.startOfWeek)
		assert.Equal(t, test.expectedYear, year, "Year should be equal")
		assert.Equal(t, test.expectedWeek, week, "Week should be equal")
	}
}

// Test for GetWeekStartEndByStartOfWeek
func TestGetWeekStartEndByStartOfWeek(t *testing.T) {
	tests := []struct {
		date          time.Time
		startOfWeek   time.Weekday
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Monday, time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 2, 23, 59, 59, 999999999, time.UTC)},
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Sunday, time.Date(2024, 5, 26, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 23, 59, 59, 999999999, time.UTC)},
	}

	for _, test := range tests {
		start, end := GetWeekStartEndByStartOfWeek(test.date, test.startOfWeek)
		assert.True(t, start.Equal(test.expectedStart), "Start date should be equal")
		assert.True(t, end.Equal(test.expectedEnd), "End date should be equal")
	}
}

// Test for GetDateOfWeekdayByStartOfWeek
func TestGetDateOfWeekdayByStartOfWeek(t *testing.T) {
	tests := []struct {
		date           time.Time
		targetWeekday  time.Weekday
		startOfWeek    time.Weekday
		expectedTarget time.Time
	}{
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Wednesday, time.Monday, time.Date(2024, 5, 29, 0, 0, 0, 0, time.UTC)},
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Sunday, time.Sunday, time.Date(2024, 5, 26, 0, 0, 0, 0, time.UTC)},
		{time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC), time.Friday, time.Sunday, time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		targetDate := GetDateOfWeekdayByStartOfWeek(test.date, test.targetWeekday, test.startOfWeek)
		assert.True(t, targetDate.Equal(test.expectedTarget), "Target date should be equal")
	}
}
