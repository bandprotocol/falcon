package wallet

type Signer interface {
	ExportPrivateKey() (string, error)
	GetName() string
	GetAddress() (addr string)
	Sign(data []byte) ([]byte, error)
}

type Wallet interface {
	SavePrivateKey(name string, privKey string) (addr string, err error)
	SaveRemoteSignerKey(name, addr, url string, key *string) error
	DeleteKey(name string) error
	GetSigners() []Signer
	GetSigner(name string) (Signer, bool)
}
