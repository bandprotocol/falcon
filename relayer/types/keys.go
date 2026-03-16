package types

// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	Mnemonic string `json:"mnemonic,omitempty"`
	Address  string `json:"address"`
	KeyName  string `json:"key_name,omitempty"`
}

// NewKeyOutput creates a new instance of KeyOutput
func NewKeyOutput(mnemonic string, address string, keyName string) *KeyOutput {
	return &KeyOutput{
		Mnemonic: mnemonic,
		Address:  address,
		KeyName:  keyName,
	}
}
