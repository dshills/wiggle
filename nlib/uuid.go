package nlib

import (
	"crypto/rand"
	"fmt"
	"io"
)

func generateUUID() string {
	// Create a slice to hold 16 bytes (128 bits) for the UUID
	u := make([]byte, 16)

	// Read 16 random bytes into the slice
	_, err := io.ReadFull(rand.Reader, u)
	if err != nil {
		panic(err) // You could handle this more gracefully
	}

	// Set the version to 4 (UUIDv4)
	u[6] = (u[6] & 0x0f) | 0x40

	// Set the variant to RFC 4122 (which is 8, 9, A, or B in hex)
	u[8] = (u[8] & 0x3f) | 0x80

	// Format as UUID string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:])
}
