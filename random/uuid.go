package random

import (
	"crypto/rand"
	"fmt"
)

const UUIDFormat = "%08x-%04x-%04x-%04x-%012x"

// UUID creates a new, version 4 uuid
func UUID() (string, error) {
	// UUID representation compliant with specification described in RFC 4122.
	var u [16]byte

	if _, err := rand.Read(u[:]); err != nil {
		return "", err
	}

	// SetVersion sets version bits.
	u[6] = (u[6] & 0x0f) | (4 << 4)
	// SetVariant sets variant bits as described in RFC 4122.
	u[8] = (u[8] & 0xbf) | 0x80

	return fmt.Sprintf(UUIDFormat, u[:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
