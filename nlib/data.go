package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// Ensure that StringData implements the node.DataCarrier interface
var _ node.DataCarrier = (*StringData)(nil)

// StringData is a simple implementation of the node.DataCarrier interface.
// It holds string data and provides methods to retrieve the data in various formats (string, JSON).
type StringData struct {
	data string // The string data held by the StringData instance
}

// NewStringData creates a new instance of StringData with the provided string data.
// This is a constructor function that returns a pointer to the newly created StringData.
func NewStringData(data string) *StringData {
	return &StringData{data: data}
}

// Vector returns nil because StringData does not support vectorized data.
// This method is part of the node.DataCarrier interface, which requires a Vector method.
func (d *StringData) Vector() []float32 {
	return nil
}

// String returns the string representation of the data held by StringData.
// This method implements the String method required by node.DataCarrier and allows
// retrieval of the raw string data.
func (d *StringData) String() string {
	return d.data
}

// JSON returns the data as a JSON-formatted byte slice. The data is embedded
// in a basic JSON structure as a string. This method satisfies the JSON method of
// the node.DataCarrier interface, which expects a JSON-compatible output.
func (d *StringData) JSON() []byte {
	return []byte(fmt.Sprintf("{ %s }", d.data))
}
