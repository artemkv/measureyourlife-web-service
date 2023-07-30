package app

import (
	"time"
)

const (
	METRIC_VALUES_MAX_TOTAL = 200
)

func isUserIdValid(userId string) bool {
	return userId != ""
}

func isEmailValid(email string) bool {
	// TODO: check email format
	return email != ""
}

func isDateValid(date string) bool {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return false
	}

	if d.Year() < 1900 || d.Year() > 2100 {
		return false
	}

	return true
}

func isMetricValueListLengthValid(dayStats dayStatsData) bool {
	return len(dayStats.MetricValues) <= METRIC_VALUES_MAX_TOTAL
}
