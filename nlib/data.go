package nlib

import (
	"github.com/dshills/wiggle/node"
)

// Ensure that StringData implements the node.DataCarrier interface
var _ node.DataCarrier = (*Carrier)(nil)

type Carrier struct {
	TextData   string
	JSONData   []byte
	VectorData [][]float32
	URLData    []string
	Base64Data []string
}

func NewTextCarrier(txt string) *Carrier {
	return &Carrier{TextData: txt}
}

func NewVectorCarrier(vec [][]float32) *Carrier {
	return &Carrier{VectorData: vec}
}

func (c *Carrier) Vector() [][]float32 {
	return c.VectorData
}

func (c *Carrier) String() string {
	return c.TextData
}

func (c *Carrier) JSON() []byte {
	return c.JSONData
}

func (c *Carrier) Base64() []string {
	return c.Base64Data
}

func (c *Carrier) ImageURLs() []string {
	return c.URLData
}
