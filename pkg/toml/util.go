package toml

import (
	"os"
	"fmt"
	"github.com/naoina/toml"
	"github.com/SeerUK/tid/pkg/types"
)

// BoltDatabaseFilename is the name of the database file name on disk.
const TomlConfigFilename = "config.toml"

// Open opens the configuration file or it creates it if it doesn't exist already.
func Open(tidDir string) (types.TomlConfig, error) {
	config := types.NewTomlConfig()

	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return config, err
	}

	var configFilePath = fmt.Sprintf("%s/%s", tidDir, TomlConfigFilename)

	f, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	fileStat, err := f.Stat();
	if err != nil {
		panic(err)
	}

	if fileStat.Size() == 0 {
		return config, nil
	}

	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config, nil
}

