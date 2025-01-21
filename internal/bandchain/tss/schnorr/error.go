// Copyright (c) 2014 Conformal Systems LLC.
// Copyright (c) 2015-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package schnorr

// ErrorKind identifies a kind of error.  It has full support for errors.Is
// and errors.As, so the caller can directly check against an error kind
// when determining the reason for an error.
type ErrorKind string

// These constants are used to identify a specific RuleError.
const (
	// ErrSigTooShort is returned when a signature that should be a Schnorr
	// signature is too short.
	ErrSigTooShort = ErrorKind("ErrSigTooShort")

	// ErrSigTooLong is returned when a signature that should be a Schnorr
	// signature is too long.
	ErrSigTooLong = ErrorKind("ErrSigTooLong")

	// ErrSigRTooBig is returned when a signature has r with a value that is
	// greater than or equal to the prime of the field underlying the group.
	ErrSigRTooBig = ErrorKind("ErrSigRTooBig")

	// ErrSigSTooBig is returned when a signature has s with a value that is
	// greater than or equal to the group order.
	ErrSigSTooBig = ErrorKind("ErrSigSTooBig")

	// ErrSigA1TooBig is returned when a complaint signature has a1 with a value that is
	// greater than or equal to the prime of the field underlying the group.
	ErrSigA1TooBig = ErrorKind("ErrSigA1TooBig")

	// ErrSigA2TooBig is returned when a complaint signature has a2 with a value that is
	// greater than or equal to the prime of the field underlying the group.
	ErrSigA2TooBig = ErrorKind("ErrSigA2TooBig")

	// ErrSigZTooBig is returned when a complaint signature has s with a value that is
	// greater than or equal to the group order.
	ErrSigZTooBig = ErrorKind("ErrSigZTooBig")

	// ErrPrivateKeyZero is returned when a private key is zero.
	ErrPrivateKeyZero = ErrorKind("ErrPrivateKeyZero")

	// ErrNotOnCurve is returned when a public key is not on curve.
	ErrNotOnCurve = ErrorKind("ErrNotOnCurve")

	// ErrRInfinity is returned when a calculated R is at infinity.
	ErrRInfinity = ErrorKind("ErrRInfinity")

	// ErrIncorrectR is returned when a calculated R is not given R.
	ErrIncorrectR = ErrorKind("ErrIncorrectR")
)

// Error satisfies the error interface and prints human-readable errors.
func (e ErrorKind) Error() string {
	return string(e)
}

// Error identifies an error related to a schnorr signature. It has full
// support for errors.Is and errors.As, so the caller can ascertain the
// specific reason for the error by checking the underlying error.
type Error struct {
	Err         error
	Description string
}

// Error satisfies the error interface and prints human-readable errors.
func (e Error) Error() string {
	return e.Description
}

// Unwrap returns the underlying wrapped error.
func (e Error) Unwrap() error {
	return e.Err
}

// signatureError creates an Error given a set of arguments.
func signatureError(kind ErrorKind, desc string) Error {
	return Error{Err: kind, Description: desc}
}
