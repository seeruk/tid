package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// The xtime package contains types and functions that provide supplementary functionality to the
// built in standard library 'time' package. It also provides some additional time-related
// functionality and types that help things like output.

// DurationFormat is an "enum" of the different display time formats. A DurationFormat can be used in
// conjunction with FormatDuration to control how some output duration is displayed.
type DurationFormat int

// All possible duration formats.
const (
	FormatDecimal DurationFormat = iota
	FormatText
)

var formats = map[string]DurationFormat{
	"decimal": FormatDecimal,
	"text":    FormatText,
}

// UnmarshalTOML takes a raw TOML time format value and attempts to parse the value as a
// DurationFormat. The value passed to this method should be a byte array of a quoted string (i.e.
// the raw TOML value), the method will remove the quotes.
func (f *DurationFormat) UnmarshalTOML(bytes []byte) error {
	text, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}

	text = strings.ToLower(text)

	format, ok := formats[text]
	if !ok {
		return fmt.Errorf("xtime: Invalid DurationFormat '%s'", text)
	}

	*f = format

	return nil
}

// Weekday is a type used to extend the built-in 'time.Weekday' type to add new methods to help with
// things like parsing weekday strings into a stricter type.
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

// UnmarshalTOML takes a raw TOML weekday string value and attempts to parse the value as a Weekday.
// The value passed to this method should be a byte array of a quoted string (i.e. the raw TOML
// value), the method will remove the quotes.
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

// LastWeekday finds the date of the most recent occurrence of a given weekday in the past.
func LastWeekday(weekday time.Weekday) time.Time {
	date := Date(time.Now())

	for date.Weekday() != weekday {
		date = date.AddDate(0, 0, -1)
	}

	return date
}

// FormatDuration returns the given time.Duration as a string in the given DurationFormat.
func FormatDuration(duration time.Duration, timeFormat DurationFormat) string {
	switch timeFormat {
	case FormatDecimal:
		return strconv.FormatFloat(duration.Hours(), 'f', 2, 64)
	case FormatText:
		fallthrough
	default:
		return duration.String()
	}
}
