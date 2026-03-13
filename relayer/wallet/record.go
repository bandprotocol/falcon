package wallet

import (
	"fmt"
	"path"
	"strings"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

const (
	LocalSignerType  = "local"
	RemoteSignerType = "remote"
)

// KeyRecord stores signer info on disk. All chain types use this unified format.
type KeyRecord struct {
	Type    string  `toml:"type"`
	Address string  `toml:"address,omitempty"`
	URL     string  `toml:"url,omitempty"`
	Key     *string `toml:"key,omitempty"`
}

// NewKeyRecord creates a new KeyRecord.
func NewKeyRecord(signerType, address, url string, key *string) KeyRecord {
	return KeyRecord{
		Type:    signerType,
		Address: address,
		URL:     url,
		Key:     key,
	}
}

// LoadKeyRecords loads all files in `dir/*.toml` into KeyRecord.
func LoadKeyRecords(dir string) (map[string]KeyRecord, error) {
	filePaths, err := os.ListFilePaths(dir)
	if err != nil {
		return nil, err
	}

	records := make(map[string]KeyRecord)
	for _, filePath := range filePaths {
		b, err := os.ReadFileIfExist(filePath)
		if err != nil {
			return nil, err
		}

		var record KeyRecord
		if err := toml.Unmarshal(b, &record); err != nil {
			return nil, err
		}

		name, err := ExtractKeyName(filePath)
		if err != nil {
			return nil, err
		}
		records[name] = record
	}

	return records, nil
}

// ExtractKeyName returns the filename (the key name) without its extension, or an error if empty.
func ExtractKeyName(filePath string) (string, error) {
	fileName := path.Base(filePath)

	keyName := strings.TrimSuffix(fileName, path.Ext(fileName))
	if keyName == "" {
		return "", fmt.Errorf("wrong keyname format")
	}

	return keyName, nil
}
