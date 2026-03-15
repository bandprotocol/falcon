package wallet

type Signer interface {
	ExportPrivateKey() (string, error)
	GetName() string
	GetAddress() string
	Sign(payload []byte, tssPayload TssPayload) ([]byte, error)
}

type Wallet interface {
	SaveByPrivateKey(name string, privateKey string) (addr string, err error)
	SaveByMnemonic(name string, mnemonic string, coinType uint32, account uint, index uint) (addr string, err error)
	SaveRemoteSignerKey(name, addr, url string, key string) error
	DeleteKey(name string) error
	GetSigners() []Signer
	GetSigner(name string) (Signer, bool)
}
