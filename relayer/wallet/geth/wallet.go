package geth

import (
	"crypto/ecdsa"
	"fmt"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Wallet = &GethWallet{}

const (
	LocalSignerType  = "local"
	RemoteSignerType = "remote"
)

// GethWallet manages local and remote signers for a specific chain.
type GethWallet struct {
	Passphrase string
	Store      *keystore.KeyStore
	Signers    map[string]wallet.Signer
	HomePath   string
	ChainName  string
}

// NewGethWallet creates a new NewGethWallet instance.
func NewGethWallet(passphrase, homePath, chainName string) (*GethWallet, error) {
	// create keystore
	keyStoreDir := path.Join(getEVMKeyStoreDir(homePath, chainName)...)
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// load signer records
	signerRecordDir := path.Join(getEVMSignerRecordDir(homePath, chainName)...)
	signerRecords, err := LoadSignerRecord(signerRecordDir)
	if err != nil {
		return nil, err
	}

	// load remote signer records
	remoteSignerDir := path.Join(getEVMRemoteSignerRecordDir(homePath, chainName)...)
	remoteSignerRecords, err := LoadRemoteSignerRecord(remoteSignerDir)
	if err != nil {
		return nil, err
	}

	signers := make(map[string]wallet.Signer)
	for name, signerRecord := range signerRecords {
		gethAddr, err := HexToETHAddress(signerRecord.Address)
		if err != nil {
			return nil, err
		}

		var signer wallet.Signer
		switch signerType := signerRecord.Type; signerType {
		case LocalSignerType:
			accs, err := store.Find(accounts.Account{Address: gethAddr})
			if err != nil {
				return nil, err
			}

			// need to export the key due to no direct access to the private key
			b, err := store.Export(accs, passphrase, passphrase)
			if err != nil {
				return nil, err
			}

			gethKey, err := keystore.DecryptKey(b, passphrase)
			if err != nil {
				return nil, err
			}

			signer = NewLocalSigner(name, gethKey.PrivateKey)
		case RemoteSignerType:
			remoteSignerRecord, ok := remoteSignerRecords[name]
			if !ok {
				return nil, fmt.Errorf("no remote signer record found: %s", name)
			}

			signer, err = NewRemoteSigner(name, gethAddr, remoteSignerRecord.Url)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf(
				"unsupported signer type %s for chain %s, key %s",
				signerRecord.Type,
				chainName,
				name,
			)
		}

		signers[name] = signer
	}

	return &GethWallet{
		Passphrase: passphrase,
		Store:      store,
		Signers:    signers,
		HomePath:   homePath,
		ChainName:  chainName,
	}, nil
}

// SavePrivateKey imports the ECDSA key into the keystore and writes its signer record.
func (w *GethWallet) SavePrivateKey(name string, privKey *ecdsa.PrivateKey) (addr string, err error) {
	// check if the key name exists
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	// derive the Ethereum address from the pubkey and check exist or not
	addr = crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	// check if the address is already added
	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	// save the signer
	_, err = w.Store.ImportECDSA(privKey, w.Passphrase)
	if err != nil {
		return "", err
	}

	if err := w.saveSignerRecord(name, SignerRecord{
		Address: addr,
		Type:    LocalSignerType,
	}); err != nil {
		return "", err
	}

	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer under the given name,
// storing its address and service URL as on‐disk records.
func (w *GethWallet) SaveRemoteSignerKey(name, address, url string) error {
	// check if the key name exists
	if _, ok := w.Signers[name]; ok {
		return fmt.Errorf("key name exists: %s", name)
	}

	// validate address
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}

	// check if the address is already added
	if w.IsAddressExist(address) {
		return fmt.Errorf("address exists: %s", name)
	}

	// save the signer
	if err := w.saveSignerRecord(name, SignerRecord{
		Address: address,
		Type:    RemoteSignerType,
	}); err != nil {
		return err
	}
	if err := w.saveRemoteSignerRecord(name, RemoteSignerRecord{
		Url: url,
	}); err != nil {
		return err
	}

	return nil
}

// DeleteKey removes the signer named name, deleting its keystore record
// if local, or its record files if remote.
func (w *GethWallet) DeleteKey(name string) error {
	// check if the key name exists
	signer, ok := w.Signers[name]
	if !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	switch signer.(type) {
	case *LocalSigner:
		addr := common.HexToAddress(signer.GetAddress())
		if err := w.Store.Delete(accounts.Account{Address: addr}, w.Passphrase); err != nil {
			return err
		}
	case *RemoteSigner:
		if err := w.deleteRemoteSignerRecord(name); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown signer type for key %q", name)
	}

	if err := w.deleteSignerRecord(name); err != nil {
		return err
	}

	return nil
}

// GetSigners lists all signers.
func (w *GethWallet) GetSigners() []wallet.Signer {
	signers := make([]wallet.Signer, 0, len(w.Signers))
	for _, signer := range w.Signers {
		signers = append(signers, signer)
	}

	return signers
}

// GetSigner returns the signer with the given name and a flag indicating if it was found.
func (w *GethWallet) GetSigner(name string) (wallet.Signer, bool) {
	signer, ok := w.Signers[name]
	return signer, ok
}

// IsAddressExist returns true if the given address is already added.
func (w *GethWallet) IsAddressExist(address string) bool {
	for _, signer := range w.Signers {
		if common.HexToAddress(signer.GetAddress()).Cmp(common.HexToAddress(address)) == 0 {
			return true
		}
	}
	return false
}

// saveSignerRecord writes the SignerRecord map to the file
func (w *GethWallet) saveSignerRecord(name string, signerRecord SignerRecord) error {
	b, err := toml.Marshal(signerRecord)
	if err != nil {
		return err
	}

	return os.Write(b, append(getEVMSignerRecordDir(w.HomePath, w.ChainName), fmt.Sprintf("%s.toml", name)))
}

// saveRemoteSignerRecord writes the RemoteSignerRecord map to the file
func (w *GethWallet) saveRemoteSignerRecord(name string, remoteSignerRecord RemoteSignerRecord) error {
	b, err := toml.Marshal(remoteSignerRecord)
	if err != nil {
		return err
	}

	return os.Write(b, append(getEVMRemoteSignerRecordDir(w.HomePath, w.ChainName), fmt.Sprintf("%s.toml", name)))
}

// deleteSignerRecord deletes the SignerRecord file
func (w *GethWallet) deleteSignerRecord(name string) error {
	dir := path.Join(getEVMSignerRecordDir(w.HomePath, w.ChainName)...)
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	return os.DeletePath(filePath)
}

// deleteRemoteSignerRecord deletes the RemoteSignerRecord file
func (w *GethWallet) deleteRemoteSignerRecord(name string) error {
	dir := path.Join(getEVMRemoteSignerRecordDir(w.HomePath, w.ChainName)...)
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	return os.DeletePath(filePath)
}
