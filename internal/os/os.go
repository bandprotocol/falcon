package os

import (
	"os"
	"path"
)

// CheckAndCreateFolder checks if the folder exists and creates it if it doesn't.
func CheckAndCreateFolder(path string) error {
	if exist, err := IsPathExist(path); err != nil {
		return err
	} else if !exist {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

// IsPathExist checks if the path exists.
func IsPathExist(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// ReadFileIfExist reads file from the given path. It returns empty data with no error,
// if the file doesn't exist. Otherwise, it returns an error.
func ReadFileIfExist(path string) ([]byte, error) {
	if exist, err := IsPathExist(path); err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}

	return os.ReadFile(path)
}

// Write writes the given data to the file at the given path. It also creates the folders if they don't exist.
func Write(data []byte, paths []string) error {
	// Create folders if they don't exist
	folderPath := ""
	for _, p := range paths[:len(paths)-1] {
		folderPath = path.Join(folderPath, p)
		if err := CheckAndCreateFolder(folderPath); err != nil {
			return err
		}
	}

	// Write the data to the file
	filePath := path.Join(folderPath, paths[len(paths)-1])
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}
