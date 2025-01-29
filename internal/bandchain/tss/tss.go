package tss

// GroupID represents the ID of a group.
type GroupID uint64

// MemberID represents the ID of a member.
// Please note that the MemberID can only be 1, 2, 3, ..., 2**64 - 1
type MemberID uint64

// SigningID represents the ID of a signing.
type SigningID uint64

// Scalar represents a scalar value stored as bytes.
// It uses secp256k1.ModNScalar and secp256k1.PrivateKey as a base implementation for serialization and parsing.
type Scalar []byte

// Point represents a point (x, y, z) stored as bytes.
// It uses secp256k1.JacobianPoint and secp256k1.PublicKey as base implementations for serialization and parsing.
type Point []byte

// Signature represents a signature (r, s) stored as bytes.
// It uses schnorr.Signature as a base implementation for serialization and parsing.
type Signature []byte
