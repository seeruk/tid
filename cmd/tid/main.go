package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/SeerUK/tid/pkg/cli"
	"github.com/SeerUK/tid/pkg/state/bolt"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/tracking"

	boltdb "github.com/boltdb/bolt"
)

func main() {
	db := getBoltDB()
	defer db.Close()

	backend := bolt.NewBoltBackend(db)

	// Initialise the backend, preparing it for use, ensuring it's up-to-date.
	migrate.Backend(backend)

	factory := tracking.NewStandardFactory(backend)
	kernel := cli.NewTidKernel(backend, factory)

	application := cli.CreateApplication()
	application.AddCommands(cli.GetCommands(kernel))

	os.Exit(application.Run(os.Args[1:]))
}

// getBoltDB gets a Bolt DB instance.
func getBoltDB() *boltdb.DB {
	// Open Bolt database.
	db, err := bolt.Open(lookupTidDir())
	fatal(err)

	return db
}

// fatal kills the application upon error.
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// lookupTidDir returns the location to store all tid files.
func lookupTidDir() string {
	// @todo: Does this belong in here?

	usr, err := user.Current()
	fatal(err)

	return fmt.Sprintf("%s/.tid", usr.HomeDir)
}
