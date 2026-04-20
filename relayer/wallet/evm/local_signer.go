package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*LocalSigner)(nil)

// LocalSigner is signer that locally stored ECDSA private key.
type LocalSigner struct {
	Name       string
	privateKey *ecdsa.PrivateKey
}

// NewLocalSigner creates a new LocalSigner instance
func NewLocalSigner(
	name string,
	privateKey *ecdsa.PrivateKey,
) *LocalSigner {
	return &LocalSigner{
		Name:       name,
		privateKey: privateKey,
	}
}

// ExportPrivateKey returns the signer's private key.
func (l *LocalSigner) ExportPrivateKey() (string, error) {
	b := crypto.FromECDSA(l.privateKey)

	return hex.EncodeToString(b), nil
}

// GetName returns the signer's key name.
func (l *LocalSigner) GetName() (addr string) {
	return l.Name
}

// GetAddress returns the signer's address.
func (l *LocalSigner) GetAddress() string {
	return crypto.PubkeyToAddress(l.privateKey.PublicKey).String()
}

// sign Keccak256-hashes data and signs it with the local private key, returning the raw signature.
func (l *LocalSigner) sign(data []byte) ([]byte, error) {
	hash := crypto.Keccak256(data)
	return crypto.Sign(hash, l.privateKey)
}

// Sign unmarshals payload as a SignerPayload JSON, builds and signs the transaction locally,
// and returns the EIP-2718 encoded signed transaction bytes.
func (l *LocalSigner) Sign(payload []byte, _ wallet.TssPayload) ([]byte, error) {
	var sp SignerPayload
	if err := json.Unmarshal(payload, &sp); err != nil {
		return nil, err
	}

	chainID := new(big.Int).SetUint64(sp.ChainID)
	to := gethcommon.HexToAddress(sp.To)
	value := decimal.NewFromInt(0).BigInt()

	var tx *gethtypes.Transaction
	switch {
	case len(sp.GasPrice) > 0:
		tx = gethtypes.NewTx(&gethtypes.LegacyTx{
			Nonce:    sp.Nonce,
			To:       &to,
			Value:    value,
			Gas:      sp.GasLimit,
			GasPrice: new(big.Int).SetBytes(sp.GasPrice),
			Data:     sp.Data,
		})
	case len(sp.GasFeeCap) > 0:
		tx = gethtypes.NewTx(&gethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     sp.Nonce,
			To:        &to,
			Value:     value,
			Gas:       sp.GasLimit,
			GasFeeCap: new(big.Int).SetBytes(sp.GasFeeCap),
			GasTipCap: new(big.Int).SetBytes(sp.GasTipCap),
			Data:      sp.Data,
		})
	default:
		return nil, fmt.Errorf("signer payload has neither GasPrice nor GasFeeCap")
	}

	signedTx, err := l.SignTx(tx, chainID)
	if err != nil {
		return nil, err
	}

	return signedTx.MarshalBinary()
}

// SignTx RLP-encodes the transaction, signs it with the local private key, and returns the signed transaction.
func (l *LocalSigner) SignTx(tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	var (
		rlpEncoded []byte
		err        error
		gethSigner gethtypes.Signer
	)

	switch tx.Type() {
	case gethtypes.LegacyTxType:
		rlpEncoded, err = rlp.EncodeToBytes(
			[]interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				chainID, uint(0), uint(0),
			},
		)
		if err != nil {
			return nil, err
		}

		gethSigner = gethtypes.NewEIP155Signer(chainID)
	case gethtypes.DynamicFeeTxType:
		rlpEncoded, err = rlp.EncodeToBytes(
			[]interface{}{
				chainID,
				tx.Nonce(),
				tx.GasTipCap(),
				tx.GasFeeCap(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.AccessList(),
			},
		)
		if err != nil {
			return nil, err
		}

		rlpEncoded = append([]byte{tx.Type()}, rlpEncoded...)
		gethSigner = gethtypes.NewLondonSigner(chainID)
	default:
		return nil, fmt.Errorf("unsupported tx type: %d", tx.Type())
	}

	signature, err := l.sign(rlpEncoded)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(gethSigner, signature)
}
