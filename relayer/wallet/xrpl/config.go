package xrpl

import (
	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

// KeyRecord stores XRPL signer info on disk.
type KeyRecord struct {
	Type    string  `toml:"type"`
	Address string  `toml:"address,omitempty"`
	Url     string  `toml:"url,omitempty"`
	Key     *string `toml:"key,omitempty"`
}

// NewKeyRecord creates a new KeyRecord.
func NewKeyRecord(signerType, address, url string, key *string) KeyRecord {
	return KeyRecord{
		Type:    signerType,
		Address: address,
		Url:     url,
		Key:     key,
	}
}

// LoadKeyRecord loads all files in `path/*.toml` into KeyRecord.
func LoadKeyRecord(path string) (map[string]KeyRecord, error) {
	filePaths, err := os.ListFilePaths(path)
	if err != nil {
		return nil, err
	}

	keyRecords := make(map[string]KeyRecord)
	for _, filePath := range filePaths {
		b, err := os.ReadFileIfExist(filePath)
		if err != nil {
			return nil, err
		}

		var keyRecord KeyRecord
		if err := toml.Unmarshal(b, &keyRecord); err != nil {
			return nil, err
		}

		name, err := ExtractKeyName(filePath)
		if err != nil {
			return nil, err
		}
		keyRecords[name] = keyRecord
	}

	return keyRecords, nil
}
