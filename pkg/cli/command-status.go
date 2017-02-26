package cli

import (
	"regexp"
	"strings"
	"time"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

func StatusCommand(gateway tracking.Gateway) console.Command {
	var short bool

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewBoolValue(&short),
			"-s, --short",
			"Show shortened output?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		if !status.IsActive() {
			output.Println("status: There is no active timer running")
			return nil
		}

		entryRef := status.TimeSheetEntry()

		date, err := time.Parse(tracking.KeyTimeSheetFmt, entryRef.Date)
		if err != nil {
			return err
		}

		sheet, err := gateway.FindTimeSheet(date)
		if err != nil {
			return err
		}

		entry := sheet.Message.Entries[entryRef.Index]

		// Update the duration, we're not persisting it in this command though.
		sheet.UpdateEntryDuration(status)

		if short {
			duration := formatDuration(timeSinceStartTime(entry.StartTime))

			output.Printf("%s on %s\n", duration, entry.Note)
		} else {
			// @todo: Should also show start time, and maybe when pausing is in, also the different
			// start times that there have been (or at least a friendly way of showing that, like
			// a timeline type thing?)
			output.Println("@todo: Long")
		}

		return nil
	}

	return console.Command{
		Name:        "status",
		Description: "View the current status. What are you tracking?",
		Configure:   configure,
		Execute:     execute,
	}
}

// @todo: Move me somewhere useful (timex? this is linked closer to an entry though)
func timeSinceStartTime(startTime uint64) time.Duration {
	startUnix := time.Unix(int64(startTime), 0)

	return time.Since(startUnix)
}

// @todo: Move me somewhere useful (timex?)
func formatDuration(duration time.Duration) string {
	toFormat := duration.String()

	// Grab the fractions of seconds in a group.
	reg := regexp.MustCompile(`\d{1,2}(\.\d+)s`)
	matches := reg.FindStringSubmatch(toFormat)

	// Replace the fractions in the seconds with nothing. This is done separately to the rest of the
	// duration to ensure accuracy.
	seconds := strings.Replace(matches[0], matches[1], "", -1)

	return strings.Replace(toFormat, matches[0], seconds, -1)
}
