// Package uuid provides the UUID type and
// helper functions.
package uuid

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// seed must never change
const seed = "bsa9phh6keet1ogh9ChoeNoK1jae8ro0"

// dash represents - as a byte
const dash byte = '-'

// UUID is a type used as identifier, compliant with
// Version 5 of RFC 4122:
// https://www.ietf.org/rfc/rfc4122.txt
type UUID string

// New creates a new UUID given a string.
func New(input string) UUID {
	id := hmac.New(sha1.New, []byte(seed))
	// we ignore errors, since it is not
	// testable
	_, _ = id.Write([]byte(input))
	b := id.Sum(nil)
	return UUID(Canonical(b))
}

// String returns a string representation of the
// UUID, implementing the Stringer interface.
func (u UUID) String() string {
	return string(u)
}

// Canonical converts an array of bytes into the
// canonical representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func Canonical(b []byte) string {
	if len(b) < 16 {
		z := make([]byte, 16-len(b))
		b = append(z, b...)
	}

	buf := make([]byte, 36)

	hex.Encode(buf[0:8], b[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], b[10:16])

	return string(buf)
}
