package nlib

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/dshills/wiggle/node"
)

// ConnectChain connects a series of nodes in sequence. It starts with the first node
// and connects each subsequent node in the provided chain. This method ensures that the nodes are
// connected linearly in the order they appear in the chain.
func ConnectChain(first node.Node, chain ...node.Node) {
	curNode := first
	for _, n := range chain {
		curNode.Connect(n) // Connect the current node to the next one in the chain
		curNode = n        // Move to the next node
	}
}

// GenerateUUID generates a UUID (Universally Unique Identifier) using the UUIDv4 format.
// It creates 16 random bytes, sets the appropriate version and variant bits, and returns
// the UUID as a string in the standard format. If there's an error generating the random bytes,
// it returns an empty string and the error.
func GenerateUUID() (string, error) {
	// Create a slice to hold 16 bytes (128 bits) for the UUID
	u := make([]byte, 16)
	// Read 16 random bytes into the slice
	_, err := io.ReadFull(rand.Reader, u)
	if err != nil {
		return "", err // Return the error if random bytes cannot be generated
	}
	// Set the version to 4 (UUIDv4) by manipulating the appropriate byte
	u[6] = (u[6] & 0x0f) | 0x40
	// Set the variant to RFC 4122 (8, 9, A, or B in hex) by manipulating the appropriate byte
	u[8] = (u[8] & 0x3f) | 0x80
	// Format the 16-byte array as a UUID string and return it
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}

// FilterMetaKey filters the metadata of a signal to return only the metadata that matches the given key.
// It iterates over the signal's metadata, checks if the key matches, and appends the matching metadata to the result slice.
// It returns the filtered slice of metadata.
func FilterMetaKey(sig node.Signal, key string) []node.Meta {
	m := []node.Meta{}
	for _, meta := range sig.Meta {
		if meta.Key == key { // Check if the meta key matches the provided key
			m = append(m, meta) // Add the matching meta to the result slice
		}
	}
	return m // Return the filtered metadata
}
