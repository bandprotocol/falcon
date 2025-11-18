package geth

import (
	"crypto/ecdsa"
	"fmt"
	"path"
	"time"

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
	start := time.Now()

	// create keystore
	keystoreStart := time.Now()
	keyStoreDir := path.Join(getEVMKeyStoreDir(homePath, chainName)...)
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Printf("Keystore creation took: %v\n", time.Since(keystoreStart))

	// load signer records
	signerRecordStart := time.Now()
	signerRecordDir := path.Join(getEVMSignerRecordDir(homePath, chainName)...)
	signerRecords, err := LoadSignerRecord(signerRecordDir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Loading signer records took: %v\n", time.Since(signerRecordStart))

	// load remote signer records
	remoteSignerStart := time.Now()
	remoteSignerDir := path.Join(getEVMRemoteSignerRecordDir(homePath, chainName)...)
	remoteSignerRecords, err := LoadRemoteSignerRecord(remoteSignerDir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Loading remote signer records took: %v\n", time.Since(remoteSignerStart))

	// process signers
	processSignersStart := time.Now()
	signers := make(map[string]wallet.Signer)
	for name, signerRecord := range signerRecords {
		signerStart := time.Now()

		gethAddr, err := HexToETHAddress(signerRecord.Address)
		if err != nil {
			return nil, err
		}

		var signer wallet.Signer
		switch signerType := signerRecord.Type; signerType {
		case LocalSignerType:
			signer = NewLocalSigner(name, gethAddr, store, passphrase)
		case RemoteSignerType:
			remoteSignerRecord, ok := remoteSignerRecords[name]
			if !ok {
				return nil, fmt.Errorf("no remote signer record found: %s", name)
			}

			signer, err = NewRemoteSigner(name, gethAddr, remoteSignerRecord.Url, remoteSignerRecord.Key)
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
		fmt.Printf("Processing signer %s took: %v\n", name, time.Since(signerStart))
	}
	fmt.Printf("Processing all signers took: %v\n", time.Since(processSignersStart))

	fmt.Printf("Total NewGethWallet creation took: %v\n", time.Since(start))

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
	start := time.Now()

	// check if the key name exists
	checkKeyStart := time.Now()
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}
	fmt.Printf("Checking key name existence took: %v\n", time.Since(checkKeyStart))

	// derive the Ethereum address from the pubkey and check exist or not
	deriveAddrStart := time.Now()
	addr = crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	fmt.Printf("Deriving address took: %v\n", time.Since(deriveAddrStart))

	// check if the address is already added
	checkAddrStart := time.Now()
	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}
	fmt.Printf("Checking address existence took: %v\n", time.Since(checkAddrStart))

	// save the signer
	importKeyStart := time.Now()
	_, err = w.Store.ImportECDSA(privKey, w.Passphrase)
	if err != nil {
		return "", err
	}
	fmt.Printf("Importing ECDSA key took: %v\n", time.Since(importKeyStart))

	saveRecordStart := time.Now()
	signerRecord := NewSignerRecord(addr, LocalSignerType)
	if err := w.saveSignerRecord(name, signerRecord); err != nil {
		return "", err
	}
	fmt.Printf("Saving signer record took: %v\n", time.Since(saveRecordStart))

	fmt.Printf("Total SavePrivateKey took: %v\n", time.Since(start))

	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer under the given name,
// storing its address and service URL as on‐disk records.
func (w *GethWallet) SaveRemoteSignerKey(name, address, url string, key *string) error {
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
	signerRecord := NewSignerRecord(address, RemoteSignerType)
	if err := w.saveSignerRecord(name, signerRecord); err != nil {
		return err
	}

	remoteSignerRecord := NewRemoteSignerRecord(url, key)
	if err := w.saveRemoteSignerRecord(name, remoteSignerRecord); err != nil {
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
