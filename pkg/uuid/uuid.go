// Package uuid provides the UUID type Version 5,
// based on SHA-1 hashing (RFC 4122), and
// helper functions.
package uuid

import gouuid "github.com/satori/go.uuid"

// seed must never change
const seed = "bsa9phh6keet1ogh9ChoeNoK1jae8ro0"

// UUID is a type used as identifier, following
// https://www.ietf.org/rfc/rfc4122.txt
type UUID string

// New creates a new UUID given a string.
func New(s string) UUID {
	// we ignore errors, since we consider
	// our seed valid to be converted to UUID
	seedUUID, _ := gouuid.FromString(seed)
	u := gouuid.NewV5(seedUUID, s)
	return UUID(u.String())
}

// String returns a string representation of the
// UUID, implementing the Stringer interface.
func (uuid UUID) String() string {
	return string(uuid)
}
