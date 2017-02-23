package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/proto"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

func StartCommand(store state.Store) console.Command {
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What are you tracking?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		// @todo: Create time sheet for today, if it doesn't exist.
		// @todo: Add new entry.

		now := time.Now().Local()
		// @todo: Move format to a constant.
		key := now.Format("2006-01-02")

		output.Printf("Started tracking '%s'.\n", note)

		err := createTimeSheetIfNotExists(key)
		if err != nil {
			return err
		}

		timeSheet, err := getTimeSheet(key)
		if err != nil {
			return err
		}

		timeSheet.Entries = append(timeSheet.Entries, createEntry(note))

		// Write to the store.
		err = store.Write("", &timeSheet)
		if err != nil {
			return err
		}

		return nil
	}

	return console.Command{
		Name:        "start",
		Description: "Start a timer.",
		Configure:   configure,
		Execute:     execute,
	}
}

func createTimeSheetIfNotExists(key string) error {
	return nil
}

func getTimeSheet(key string) (proto.TimeSheet, error) {
	return proto.TimeSheet{}, nil
}

func createEntry(note string) *proto.TimeSheetEntry {
	return &proto.TimeSheetEntry{
		Note:      note,
		StartTime: uint64(time.Now().Unix()),
	}
}
