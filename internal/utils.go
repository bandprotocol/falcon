package internal

import "os"

// CheckAndCreateFolder checks if the folder exists and creates it if it doesn't.
func CheckAndCreateFolder(path string) error {
	// If the folder exists and no error, return nil
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	// If the folder does not exist, create it.
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return err
}
