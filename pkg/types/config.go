package types

type Config struct {
	Display ConfigDisplay
}

type ConfigDisplay struct {
	TimeFormat   string
	FirstWeekday string
}

// NewConfig creates a Config struct with default values
func NewConfig() Config {
	return Config{
		Display: ConfigDisplay{
			TimeFormat:   "text",
			FirstWeekday: "Monday",
		},
	}
}
