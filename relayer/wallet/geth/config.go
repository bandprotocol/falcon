package geth

import (
	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

// SignerRecord is signer information.
type SignerRecord struct {
	Address string `toml:"address"`
	Type    string `toml:"type"`
}

// NewSignerRecord creates a new SignerRecord.
func NewSignerRecord(address string, signerType string) SignerRecord {
	return SignerRecord{
		Address: address,
		Type:    signerType,
	}
}

// RemoteSignerRecord is remote signer's information.
type RemoteSignerRecord struct {
	Url string  `toml:"url"`
	Key *string `toml:"key"`
}

// NewRemoteSignerRecord creates a new RemoteSignerRecord.
func NewRemoteSignerRecord(url string, key *string) RemoteSignerRecord {
	return RemoteSignerRecord{
		Url: url,
		Key: key,
	}
}

// LoadSignerRecord loads all files in `path/*.toml` into SignerRecord.
func LoadSignerRecord(path string) (map[string]SignerRecord, error) {
	filePaths, err := os.ListFilePaths(path)
	if err != nil {
		return nil, err
	}

	signerRecords := make(map[string]SignerRecord)
	for _, filePath := range filePaths {
		b, err := os.ReadFileIfExist(filePath)
		if err != nil {
			return nil, err
		}

		// unmarshal SignerInfo
		var signerRecord SignerRecord
		err = toml.Unmarshal(b, &signerRecord)
		if err != nil {
			return nil, err
		}

		name, err := ExtractKeyName(filePath)
		if err != nil {
			return nil, err
		}
		signerRecords[name] = signerRecord
	}

	return signerRecords, nil
}

// LoadRemoteSignerRecord loads all files in `path/*.toml` into RemoteSignerRecord.
func LoadRemoteSignerRecord(path string) (map[string]RemoteSignerRecord, error) {
	filePaths, err := os.ListFilePaths(path)
	if err != nil {
		return nil, err
	}

	remoteSignerRecords := make(map[string]RemoteSignerRecord)
	for _, filePath := range filePaths {
		b, err := os.ReadFileIfExist(filePath)
		if err != nil {
			return nil, err
		}

		// unmarshal RemoterSignerRecord
		var remoteSignerRecord RemoteSignerRecord
		err = toml.Unmarshal(b, &remoteSignerRecord)
		if err != nil {
			return nil, err
		}

		name, err := ExtractKeyName(filePath)
		if err != nil {
			return nil, err
		}
		remoteSignerRecords[name] = remoteSignerRecord
	}

	return remoteSignerRecords, nil
}
