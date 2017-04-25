package types

type TomlConfig struct {
	Display TidDisplay
}

type TidDisplay struct {
	TimeFormat string
	FirstWeekDay string
}

// NewTomlConfig creates a TomlConfig struct with default values
func NewTomlConfig() TomlConfig {
	return TomlConfig{
		Display: TidDisplay{
			TimeFormat: "text",
			FirstWeekDay: "Monday",
		},
	}
}
