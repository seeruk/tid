package main

import (
	"log"
	"os"

	"github.com/SeerUK/tid/pkg/state/bolt"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/tid"
	"github.com/SeerUK/tid/pkg/tid/cli"
	"github.com/SeerUK/tid/pkg/util"

	boltdb "github.com/boltdb/bolt"

	_ "github.com/SeerUK/tid/pkg/state/migrate/versions"
	"github.com/SeerUK/tid/pkg/toml"
	"fmt"
	"github.com/SeerUK/tid/pkg/types"
)

func main() {
	dir, err := tid.GetLocalDirectory()
	fatal(err)

	tomlConfig := getTomlConfig(dir)
	// then to use the unmarshaled config...
	fmt.Println("Config name: ", tomlConfig.Owner.Name)

	db := getBoltDB(dir)
	defer db.Close()

	backend := bolt.NewBoltBackend(db)

	// Initialise the backend, preparing it for use, ensuring it's up-to-date.
	migErr := migrate.Backend(backend)
	fatal(migErr)

	factory := util.NewStandardFactory(backend)
	kernel := cli.NewTidKernel(backend, factory)

	os.Exit(cli.CreateApplication(kernel).Run(os.Args[1:], os.Environ()))
}

// getTomlConfig gets the config
func getTomlConfig(dir string) types.TomlConfig {
	tomlConfig, err := toml.Open(dir)
	fatal(err)

	return tomlConfig
}

// getBoltDB gets a Bolt DB instance.
func getBoltDB(dir string) *boltdb.DB {
	// Open Bolt database.
	db, err := bolt.Open(dir)
	fatal(err)

	return db
}

// fatal kills the application upon error.
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
