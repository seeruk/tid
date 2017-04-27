package types

import (
	"github.com/SeerUK/tid/pkg/xtime"
)

// Config represents the application configuration format.
type Config struct {
	Display ConfigDisplay
}

// ConfigDisplay represents configuration for output.
type ConfigDisplay struct {
	TimeFormat   xtime.TimeFormat
	FirstWeekday xtime.Weekday
}

// NewConfig creates a Config struct with default values.
func NewConfig() Config {
	return Config{
		Display: ConfigDisplay{
			TimeFormat:   xtime.FormatText,
			FirstWeekday: xtime.Monday,
		},
	}
}
