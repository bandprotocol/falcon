package types

// KeyOutput contains mnemonic and address of key
type Key struct {
	Mnemonic string `json:"mnemonic,omitempty"`
	Address  string `json:"address"`
	KeyName  string `json:"key_name,omitempty"`
}

// NewKeyOutput creates a new instance of KeyOutput
func NewKey(mnemonic string, address string, keyName string) *Key {
	return &Key{
		Mnemonic: mnemonic,
		Address:  address,
		KeyName:  keyName,
	}
}
