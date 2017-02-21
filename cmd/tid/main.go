package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/SeerUK/tid/pkg/cli"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/eidolon/console"
)

func main() {
	// Open Bolt database.
	bolt, err := state.OpenBolt(lookupTidDir())
	fatal(err)

	defer bolt.Close()

	// Create required Buckets.
	err = state.InitialiseBolt(bolt)
	fatal(err)

	// Pass Bolt database to create store instance.
	//store := state.NewBoltStore(bolt, state.BoltBucketTimeSheet)

	application := cli.CreateApplication()
	application.AddCommands([]console.Command{})

	os.Exit(application.Run(os.Args[1:]))
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
