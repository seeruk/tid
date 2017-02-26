package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/SeerUK/tid/pkg/cli"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

func main() {
	store := getStore()
	defer store.Close()

	tsGateway := tracking.NewGateway(store)

	application := cli.CreateApplication()
	application.AddCommands([]console.Command{
		cli.ResumeCommand(tsGateway),
		cli.StartCommand(tsGateway),
		cli.StatusCommand(tsGateway),
		cli.StopCommand(tsGateway),
	})

	os.Exit(application.Run(os.Args[1:]))
}

// getStore gets the application data store, in a ready state.
func getStore() state.Store {
	// Open Bolt database.
	bolt, err := state.OpenBolt(lookupTidDir())
	fatal(err)

	// Create required Buckets.
	err = state.InitialiseBolt(bolt)
	fatal(err)

	// Pass Bolt database to create store instance.
	return state.NewBoltStore(bolt, state.BoltBucketTimeSheet)
}

// fatal kills the application upon error.
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// lookupTidDir returns the location to store all tid files.
func lookupTidDir() string {
	usr, err := user.Current()
	fatal(err)

	return fmt.Sprintf("%s/.tid", usr.HomeDir)
}
