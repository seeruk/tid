package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/SeerUK/tid/pkg/cli"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/bolt"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"

	boltdb "github.com/boltdb/bolt"
)

func main() {
	db := getBoltDB()

	sysStore := getStore(db, bolt.BoltBucketSys)
	timeSheetStore := getStore(db, fmt.Sprintf(
		bolt.BoltBucketWorkspaceFmt,
		bolt.BoltBucketWorkspaceDefault,
	))

	defer sysStore.Close()
	defer timeSheetStore.Close()

	sysGateway := tracking.NewStoreSysGateway(sysStore)
	tsGateway := tracking.NewStoreTimesheetGateway(timeSheetStore, sysGateway)

	facade := tracking.NewFacade(sysGateway, tsGateway)

	application := cli.CreateApplication()
	application.AddCommands([]console.Command{
		cli.AddCommand(tsGateway),
		cli.EditCommand(sysGateway, tsGateway),
		cli.RemoveCommand(tsGateway, facade),
		cli.ReportCommand(sysGateway, tsGateway),
		cli.ResumeCommand(sysGateway, tsGateway),
		cli.StartCommand(sysGateway, tsGateway),
		cli.StatusCommand(sysGateway, tsGateway),
		cli.StopCommand(sysGateway, tsGateway),
		cli.WorkspaceCommand(tsGateway),
	})

	os.Exit(application.Run(os.Args[1:]))
}

// getBoltDB gets a Bolt DB instance.
func getBoltDB() *boltdb.DB {
	// Open Bolt database.
	db, err := bolt.Open(lookupTidDir())
	fatal(err)

	// Create required Buckets.
	err = bolt.Initialise(db)
	fatal(err)

	return db
}

// getStore gets the application data store, in a ready state.
func getStore(db *boltdb.DB, bucketName string) state.Store {
	return bolt.NewBoltStore(db, bucketName)
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
