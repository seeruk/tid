package tid

import (
	"fmt"
	"os/user"
)

// GetLocalDirectory returns the location to store all local tid files.
func GetLocalDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.tid", usr.HomeDir), nil
}
