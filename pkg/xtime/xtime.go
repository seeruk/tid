package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeFormat is an "enum" of the different display time formats.
type TimeFormat int

const (
	// FormatDecimal format, e.g. 1.85
	FormatDecimal TimeFormat = iota
	// FormatText format, e.g. 1h51m0s
	FormatText
)

var formats = map[string]TimeFormat{
	"decimal": FormatDecimal,
	"text":    FormatText,
}

// UnmarshalTOML takes a raw TOML time format value and attempts to parse the value as a TimeFormat.
func (f *TimeFormat) UnmarshalTOML(bytes []byte) error {
	text, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	text = strings.ToLower(text)

	format, ok := formats[text]
	if !ok {
		return fmt.Errorf("xtime: Invalid TimeFormat '%s'", text)
	}

	*f = format

	return nil
}

// A Weekday specifies a day of the week (Sunday = 0, ...).
type Weekday time.Weekday

// All possible Weekday values.
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var weekdays = map[string]Weekday{
	"sunday":    Sunday,
	"monday":    Monday,
	"tuesday":   Tuesday,
	"wednesday": Wednesday,
	"thursday":  Thursday,
	"friday":    Friday,
	"saturday":  Saturday,
}

// TimeWeekday converts this Weekday to a standard library time.Weekday.
func (w Weekday) TimeWeekday() time.Weekday {
	return time.Weekday(w)
}

// UnmarshalTOML takes a raw TOML weekday value and attempts to parse the value as a Weekday.
func (w *Weekday) UnmarshalTOML(bytes []byte) error {
	text, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	text = strings.ToLower(text)

	weekday, ok := weekdays[text]
	if !ok {
		return fmt.Errorf("xtime: Invalid Weekday '%s'", text)
	}

	*w = weekday

	return nil
}

const (
	// DateFmt is a fairly standard, date-only format for times.
	DateFmt = "2006-01-02"
)

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

// FormatDuration prints the time using the given format.
func FormatDuration(duration time.Duration, timeFormat TimeFormat) string {
	switch timeFormat {
	case FormatDecimal:
		return strconv.FormatFloat(duration.Hours(), 'f', 2, 64)
	case FormatText:
		fallthrough
	default:
		return duration.String()
	}
}
