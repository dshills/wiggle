package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// GenerateAINodeFactory will return a node.Factory function that allows an arbitrary number of AI Nodes
// to be created on-demand. Used with PartitionerNode to split workflows
func GenerateAINodeFactory(lm llm.LLM, l node.Logger, sm node.StateManager, namePrefix string) node.Factory {
	return func(count int) []node.Node {
		nodes := []node.Node{}
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("%s-%s", namePrefix, GenerateUUID())
			nodes = append(nodes, NewAINode(lm, l, sm, name))
		}
		return nodes
	}
}
