package toml

import (
	"os"
	"fmt"
	"github.com/naoina/toml"
	"github.com/SeerUK/tid/pkg/types"
)

// BoltDatabaseFilename is the name of the database file name on disk.
const TomlConfigFilename = "config.toml"

// Open opens a Bolt database, creating it if it doesn't exist already.
func Open(tidDir string) (types.TomlConfig, error) {
	var config types.TomlConfig

	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return config, err
	}

	var configFilePath = fmt.Sprintf("%s/%s", tidDir, TomlConfigFilename);

	f, err := os.Open(configFilePath)
	if err != nil {
		//f, err := os.Create(configFilePath)
		if err != nil {
			panic(err)
		} else {

		}
	}
	defer f.Close()

	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config, nil
}

