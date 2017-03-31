package timeutil

import "time"

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

	for date.Weekday() != weekday {
		date = date.AddDate(0, 0, -1)
	}

	return date
}
