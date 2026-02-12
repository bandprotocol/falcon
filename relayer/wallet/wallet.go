package wallet

type Signer interface {
	ExportPrivateKey() (string, error)
	GetName() string
	GetAddress() (addr string)
	Sign(data []byte, preSignPayload *PreSignPayload) ([]byte, error)
}

type Wallet interface {
	SaveByPrivateKey(name string, privateKey string) (addr string, err error)
	SaveByFamilySeed(name string, familySeed string) (addr string, err error)
	SaveByMnemonic(name string, mnemonic string, coinType uint32, account uint, index uint) (addr string, err error)
	SaveRemoteSignerKey(name, addr, url string, key *string) error
	DeleteKey(name string) error
	GetSigners() []Signer
	GetSigner(name string) (Signer, bool)
}

type PreSignPayload struct {
	TssMessage []byte
	RandomAddr []byte
	Signature  []byte
}

func NewPreSignPayload(tssMessage, randomAddr, signature []byte) *PreSignPayload {
	return &PreSignPayload{
		TssMessage: tssMessage,
		RandomAddr: randomAddr,
		Signature:  signature,
	}
}
