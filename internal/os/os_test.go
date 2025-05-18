package os_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	internal_os "github.com/bandprotocol/falcon/internal/os"
)

func TestCheckAndCreateFolder(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// create a folder
	err := internal_os.CheckAndCreateFolder(tmpDir)
	require.NoError(t, err)

	_, err = os.Stat(tmpDir)
	require.NoError(t, err)

	// create a folder again; shouldn't cause any error
	err = internal_os.CheckAndCreateFolder(tmpDir)
	require.NoError(t, err)
}

func TestIsPathExist(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// check if a folder exists, should return false
	exist, err := internal_os.IsPathExist(tmpDir)
	require.NoError(t, err)
	require.False(t, exist)

	err = internal_os.CheckAndCreateFolder(tmpDir)
	require.NoError(t, err)

	exist, err = internal_os.IsPathExist(tmpDir)
	require.NoError(t, err)
	require.True(t, exist)
}

func TestListFilePath(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// write two files
	err := internal_os.Write([]byte("test1"), []string{tmpDir, "test1.txt"})
	require.NoError(t, err)

	err = internal_os.Write([]byte("test2"), []string{tmpDir, "test2.txt"})
	require.NoError(t, err)

	// list and verify returned file paths
	filePaths, err := internal_os.ListFilePaths(tmpDir)
	require.NoError(t, err)
	require.Equal(t, []string{path.Join(tmpDir, "test1.txt"), path.Join(tmpDir, "test2.txt")}, filePaths)
}

func TestWrite(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// write a file
	err := internal_os.Write([]byte("test"), []string{tmpDir, "test.txt"})
	require.NoError(t, err)

	// check if the file exists
	exist, err := internal_os.IsPathExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.True(t, exist)

	// check if the file contains the correct data
	data, err := os.ReadFile(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.Equal(t, "test", string(data))

	// write a file again; shouldn't cause any error
	err = internal_os.Write([]byte("new test"), []string{tmpDir, "test.txt"})
	require.NoError(t, err)

	// check if the file contains the correct data
	data, err = os.ReadFile(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.Equal(t, "new test", string(data))
}

func TestReadFileIfExist(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// check if a file exists, should return nil
	data, err := internal_os.ReadFileIfExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.Nil(t, data)

	// write a file
	err = internal_os.Write([]byte("test"), []string{tmpDir, "test.txt"})
	require.NoError(t, err)

	// check if the file exists
	exist, err := internal_os.IsPathExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.True(t, exist)

	// check if the file contains the correct data
	data, err = internal_os.ReadFileIfExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.Equal(t, "test", string(data))

	// check if a file doesn't exist, should return nil
	data, err = internal_os.ReadFileIfExist(path.Join(tmpDir, "non-exist.txt"))
	require.NoError(t, err)
	require.Nil(t, data)
}

func TestDeletePath(t *testing.T) {
	tmpDir := path.Join(t.TempDir(), "test")

	// write a file
	err := internal_os.Write([]byte("test"), []string{tmpDir, "test.txt"})
	require.NoError(t, err)

	// check if the file exists
	exist, err := internal_os.IsPathExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.True(t, exist)

	// delete file
	err = internal_os.DeletePath(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)

	// check if a file doesn't exist, should return nil
	exist, err = internal_os.IsPathExist(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)
	require.False(t, exist)
}
