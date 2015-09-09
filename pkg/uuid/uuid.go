// package uuid provides the UUID type and
// helper functions.
package uuid

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
)

// seed must never change
const seed = "bsa9phh6keet1ogh9ChoeNoK1jae8ro0"

// UUID is a type used as identifier, following
// https://www.ietf.org/rfc/rfc4122.txt
type UUID string

// NewUUID creates a new UUID given a string.
func New(s string) UUID {
	h := hmac.New(md5.New, []byte(seed))
	// we ignore errors, since it is not
	// testable
	_, _ = h.Write([]byte(s))
	return UUID(fmt.Sprintf("%x", h.Sum(nil)))
}

// String returns a string representation of the
// UUID, implementing the Stringer interface.
func (uuid UUID) String() string {
	return string(uuid)
}
