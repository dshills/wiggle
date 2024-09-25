package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.DataCarrier = (*StringData)(nil)

type StringData struct {
	data string
}

func NewStringData(data string) *StringData {
	return &StringData{data: data}
}

func (d *StringData) Vector() []float32 {
	return nil
}

func (d *StringData) String() string {
	return d.data
}

func (d *StringData) JSON() []byte {
	return []byte(fmt.Sprintf("{ %s }", d.data))
}
