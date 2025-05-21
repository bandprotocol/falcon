package os

import (
	"fmt"
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

// ListFilePaths returns all file paths (non-dirs) immediately under dir.
// If the directory does not exist, it returns an empty slice without error.
func ListFilePaths(dirPath string) ([]string, error) {
	exist, err := IsPathExist(dirPath)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	var filePaths []string

	paths, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	for _, p := range paths {
		if p.IsDir() {
			return []string{}, fmt.Errorf("folder exists")
		}
		filePaths = append(filePaths, path.Join(dirPath, p.Name()))
	}
	return filePaths, nil
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

// DeletePath removes the file or directory at the given path.
// It returns an error if the path does not exist.
func DeletePath(path string) error {
	return os.Remove(path)
}
