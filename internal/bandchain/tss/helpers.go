package types

import "github.com/ethereum/go-ethereum/crypto"

// PaddingBytes pads a byte slice with zero bytes to a specified length.
func PaddingBytes(data []byte, length int) []byte {
	if len(data) >= length {
		return data
	}

	padding := make([]byte, length-len(data))
	return append(padding, data...)
}

// H(m)
// Hash calculates the Keccak-256 hash of the given data.
// It returns the hash value as a byte slice.
func Hash(data ...[]byte) []byte {
	return crypto.Keccak256(data...)
}
