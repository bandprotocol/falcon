package evm

import "fmt"

var (
	evmProviderErr = func(err error) error {
		return fmt.Errorf("[EVMProvider] %w", err)
	}
	evmClientErr = func(err error) error {
		return fmt.Errorf("[EVMClient] %w", err)
	}
)

var (
	// Common error
	ErrEstimateGas = func(err error, isClient bool) error {
		if !isClient {
			return evmClientErr(fmt.Errorf("failed to estimate gas: %s", err))
		}
		return evmProviderErr(fmt.Errorf("failed to estimate gas: %s", err))
	}

	// EVMProvider errors
	ErrLoadAbi = func(err error) error {
		return evmProviderErr(fmt.Errorf("failed to load abi: %w", err))
	}

	ErrUnsupportedGasType = func(gasType GasType) error {
		return evmProviderErr(fmt.Errorf("unsupported gas type: %v", gasType))
	}

	ErrPackCalldata = func(err error) error {
		return evmProviderErr(fmt.Errorf("failed to pack calldata: %s", err))
	}

	ErrQueryData = func(err error) error {
		return evmProviderErr(fmt.Errorf("failed to query data: %w", err))
	}

	ErrUnpackData = func(err error) error {
		return evmProviderErr(fmt.Errorf("failed to unpack data: %w", err))
	}

	ErrInvalidAddress = func(err error) error {
		return evmProviderErr(fmt.Errorf("invalid address: %w", err))
	}

	ErrClientConnection = func(err error) error {
		return evmProviderErr(fmt.Errorf("failed to query contract: %w", err))
	}

	ErrRelayPacketRetries = func(retries int) error {
		return evmProviderErr(fmt.Errorf("failed to relay packet after %d retries", retries))
	}

	// EVMClient errors
	ErrGetNonce = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to get nonce: %w", err))
	}

	ErrGetBlockHeight = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to get block height: %w", err))
	}

	ErrGetTxReceipt = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to get tx receipt: %w", err))
	}

	ErrGetTxByHash = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to get tx by hash: %w", err))
	}

	ErrQuery = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to query: %w", err))
	}

	ErrBroadcastTx = func(s string) error {
		return evmClientErr(fmt.Errorf("failed to broadcast tx with error: %s", s))
	}

	ErrConnectEVMChain = evmClientErr(fmt.Errorf("failed to connect to EVM chain"))

	ErrQueryBalance = func(err error) error {
		return evmClientErr(fmt.Errorf("failed to query balance: %w", err))
	}
)
