package nlib

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/dshills/wiggle/node"
)

func ConnectChain(first node.Node, chain ...node.Node) {
	curNode := first
	for _, n := range chain {
		curNode.Connect(n)
		curNode = n
	}
}

func GenerateUUID() string {
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

func NewGroup(orgID string, nodes ...node.Node) node.Group {
	grp := node.Group{
		OriginatorID: orgID,
		BatchID:      GenerateUUID(),
	}
	for _, n := range nodes {
		grp.TaskIDs = append(grp.TaskIDs, n.ID())
	}

	return grp
}

func FilterMetaKey(sig node.Signal, key string) []node.Meta {
	m := []node.Meta{}
	for _, meta := range sig.Meta {
		if meta.Key == key {
			m = append(m, meta)
		}
	}
	return m
}
