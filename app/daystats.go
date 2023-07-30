package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type dateContainerData struct {
	Date string `uri:"dt" binding:"required"`
}

type dayStatsData struct {
	MetricValues []metricValue `json:"metric_values" binding:"required"`
}

type metricValue struct {
	Id    string `json:"id" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func handleGetDayStats(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateContainer dateContainerData
	if err := c.ShouldBindUri(&dateContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isDateValid(dateContainer.Date) {
		err := fmt.Errorf("invalid value '%s' for 'date'", dateContainer.Date)
		toBadRequest(c, err)
		return
	}

	dayStats, err := getDayStats(userId, dateContainer.Date)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	if dayStats == nil {
		dayStats = &dayStatsData{
			MetricValues: []metricValue{},
		}
	}

	if dayStats.MetricValues == nil {
		dayStats.MetricValues = []metricValue{}
	}

	// TODO: this is for testing, remove when no more useful
	// time.Sleep(300 * time.Millisecond)
	/*toBadRequest(c, fmt.Errorf("Something went wrong returning day stats"))
	return*/

	toSuccess(c, dayStats)
}

func handlePostDayStats(c *gin.Context, userId string, email string) {
	// get date from URL
	var dateContainer dateContainerData
	if err := c.ShouldBindUri(&dateContainer); err != nil {
		toBadRequest(c, err)
		return
	}

	// get day stats data from the POST body
	var dayStats dayStatsData
	if err := c.ShouldBindJSON(&dayStats); err != nil {
		toBadRequest(c, err)
		return
	}

	// sanitize
	if !isDateValid(dateContainer.Date) {
		err := fmt.Errorf("invalid value '%s' for 'date'", dateContainer.Date)
		toBadRequest(c, err)
		return
	}
	if !isMetricValueListLengthValid(dayStats) {
		err := fmt.Errorf("too many metric values, max %d allowed",
			METRIC_VALUES_MAX_TOTAL)
		toBadRequest(c, err)
		return
	}

	err := updateDayStats(userId, dateContainer.Date, dayStats)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	/*toBadRequest(c, fmt.Errorf("Something went wrong saving day stats"))
	  return*/

	toSuccess(c, dayStats)
}
