package wallet

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
)

const (
	PrivKeySignerType  = "privkey"
	MnemonicSignerType = "mnemonic"
	RemoteSignerType   = "remote"
)

// KeyRecord stores signer info on disk. All chain types use this unified format.
// Secrets (private key, mnemonic+derivation params, remote API key) are stored
// in the shared keyring, never in this file.
type KeyRecord struct {
	Type    string `toml:"type"`
	Address string `toml:"address,omitempty"`
	URL     string `toml:"url,omitempty"`
}

// NewKeyRecord creates a new KeyRecord.
func NewKeyRecord(signerType, address, url string) KeyRecord {
	return KeyRecord{
		Type:    signerType,
		Address: address,
		URL:     url,
	}
}

// MnemonicSecret is the keyring payload for mnemonic-derived signers.
// Bundling the derivation parameters here keeps KeyRecord minimal and ensures
// all state needed to reconstruct a signer lives in the (encrypted) keyring.
type MnemonicSecret struct {
	Mnemonic string `json:"mnemonic"`
	CoinType uint32 `json:"coin_type"`
	Account  uint   `json:"account"`
	Index    uint   `json:"index"`
}

// EncodeMnemonicSecret serialises the mnemonic and derivation params into a
// JSON string suitable for storing in the shared keyring.
func EncodeMnemonicSecret(mnemonic string, coinType uint32, account, index uint) (string, error) {
	b, err := json.Marshal(MnemonicSecret{
		Mnemonic: mnemonic,
		CoinType: coinType,
		Account:  account,
		Index:    index,
	})
	return string(b), err
}

// DecodeMnemonicSecret deserialises a JSON string produced by EncodeMnemonicSecret.
func DecodeMnemonicSecret(secret string) (MnemonicSecret, error) {
	var ms MnemonicSecret
	err := json.Unmarshal([]byte(secret), &ms)
	return ms, err
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
