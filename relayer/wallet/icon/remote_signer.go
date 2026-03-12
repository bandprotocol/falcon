package icon

import (
	"fmt"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner signs ICON transactions via the KMS service.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new ICON RemoteSigner.
// Its signature matches wallet.RemoteSignerFactory so it can be passed
// directly to wallet.NewRemoteOnlyAdapter.
func NewRemoteSigner(name, address, url string, key *string) (wallet.Signer, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address, url, key)
	if err != nil {
		return nil, err
	}
	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign requests the remote KMS to sign the ICON transaction.
func (r *RemoteSigner) Sign(payload []byte, tss wallet.TssPayload) ([]byte, error) {
	// TODO: replace with the real proto call once SignIcon is added to fkms.proto:
	//   res, err := r.FkmsClient.SignIcon(r.ContextWithKey(), &fkmsv1.SignIconRequest{
	//       Address: r.Address, Message: payload,
	//   })
	_ = fkmsv1.ChainType_EVM // keep the import used until the real RPC exists
	return nil, fmt.Errorf("SignIcon not yet implemented")
}
