package timeutil

import (
	"time"
	"strings"
	"fmt"
)

// DateFmt is a fairly standard, date-only format for times.
const DateFmt = "2006-01-02"

// Date returns a given time, without any time on it. Only the date.
func Date(datetime time.Time) time.Time {
	date, err := time.Parse(DateFmt, datetime.Format(DateFmt))
	if err != nil {
		// This should only happen if something really crazy happens...
		panic(err)
	}

	return date
}

// LastWeekday finds the last date for the given weekday.
func LastWeekday(weekday time.Weekday) time.Time {
	date := Date(time.Now())

	for date.Weekday() !=  weekday{
		date = date.AddDate(0, 0, -1)
	}

	return date
}

// StringToWeekday transforms a string to time.Weekday
// it returns Monday if wrong day specified in config.toml
func StringToWeekday(weekday string) (time.Weekday, error){
	switch strings.ToLower(weekday) {
		case "monday":
			return time.Monday, nil
		case "tuesday":
			return time.Tuesday, nil
		case "wednesday":
			return time.Wednesday, nil
		case "thursday":
			return time.Thursday, nil
		case "friday":
			return time.Friday, nil
		case "saturday":
			return time.Saturday, nil
		case "sunday":
			return time.Sunday, nil
	}

	return time.Sunday, fmt.Errorf("Invalid first week day '%s' specifid in config.toml. Sunday is used as default", weekday)
}
