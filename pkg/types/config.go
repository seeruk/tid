package types

// Config represents the application configuration format.
type Config struct {
	Display ConfigDisplay
}

// ConfigDisplay represents configuration for output.
type ConfigDisplay struct {
	TimeFormat   string
	FirstWeekday string
}

// NewConfig creates a Config struct with default values.
func NewConfig() Config {
	return Config{
		Display: ConfigDisplay{
			TimeFormat:   "text",
			FirstWeekday: "Monday",
		},
	}
}
