package nlib

import (
	"fmt"
	"strconv"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// GenerateAINodeFactory will return a node.Factory function that allows an arbitrary number of AI Nodes
// to be created on-demand. Used with PartitionerNode to split workflows
func GenerateAINodeFactory(lm llm.LLM, mgr node.StateManager, prefix string, options node.Options) node.Factory {
	return func(count int) []node.Node {
		nodes := []node.Node{}
		for i := 0; i < count; i++ {
			suffix, err := GenerateUUID()
			if err != nil {
				suffix = strconv.Itoa(i)
			}
			options.ID = fmt.Sprintf("%s-%s", prefix, suffix)
			nodes = append(nodes, NewAINode(lm, mgr, options))
		}
		return nodes
	}
}
