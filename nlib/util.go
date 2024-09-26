package nlib

import "github.com/dshills/wiggle/node"

func ConnectChain(first node.Node, chain ...node.Node) {
	curNode := first
	for _, n := range chain {
		curNode.Connect(n)
		curNode = n
	}
}
