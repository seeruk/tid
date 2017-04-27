package toml

import (
	"os"

	"path/filepath"

	"github.com/SeerUK/tid/pkg/types"
	"github.com/naoina/toml"
)

// BoltDatabaseFilename is the name of the database file name on disk.
const TomlConfigFilename = "config.toml"

// Open opens the configuration file or it creates it if it doesn't exist already.
func Open(tidDir string) (types.Config, error) {
	config := types.NewConfig()

	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return config, err
	}

	var configFilePath = filepath.Join(tidDir, TomlConfigFilename)

	f, err := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return config, err
	}

	defer f.Close()

	fileStat, err := f.Stat()
	if err != nil {
		return config, err
	}

	if fileStat.Size() == 0 {
		return config, nil
	}

	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config, nil
}
