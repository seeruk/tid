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

	"github.com/SeerUK/tid/pkg/state/migrate"
	boltdb "github.com/boltdb/bolt"
)

func main() {
	db := getBoltDB()
	defer db.Close()

	backend := bolt.NewBoltBackend(db)

	// Initialise the backend, preparing it for use, ensuring it's up-to-date.
	migrate.Backend(backend)

	sysStore := getStore(backend, state.BackendBucketSys)
	sysGateway := tracking.NewStoreSysGateway(sysStore)

	tsStore := getStore(backend, getWorkspaceBucketName(sysGateway))
	tsGateway := tracking.NewStoreTimesheetGateway(tsStore, sysGateway)

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
		cli.WorkspaceCommand(backend, sysGateway),
	})

	os.Exit(application.Run(os.Args[1:]))
}

// getBoltDB gets a Bolt DB instance.
func getBoltDB() *boltdb.DB {
	// Open Bolt database.
	db, err := bolt.Open(lookupTidDir())
	fatal(err)

	return db
}

// getStore gets the application data store.
func getStore(db state.Backend, bucketName string) state.Store {
	return state.NewBackendStore(db, bucketName)
}

// getWorkspaceBucketName gets the name of the currently active workspace's bucket in Bolt.
func getWorkspaceBucketName(sysGateway tracking.SysGateway) string {
	// @todo: Does this belong in here?

	status, err := sysGateway.FindOrCreateStatus()
	fatal(err)

	return fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		status.Workspace,
	)
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
