package nlib

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/wiggle/node"
)

// Ensure that SimplePartitionerNode implements the node.PartitionerNode interface
var _ node.PartitionerNode = (*SimplePartitionerNode)(nil)

// SimplePartitionerNode is a node that partitions incoming signals into smaller tasks,
// processes them in parallel, and then integrates the results. It uses a partition function
// to divide the input, a node factory to create new nodes for processing, and an integration
// function to combine the results.
type SimplePartitionerNode struct {
	EmptyNode                          // Embeds the base functionality of an EmptyNode
	partitionFunc   node.PartitionerFn // Function for partitioning signal data
	integrationFunc node.IntegratorFn  // Function for integrating the partitioned results
	factory         node.Factory       // Factory function to create nodes for processing partitions
}

// NewSimplePartitionerNode creates a new SimplePartitionerNode. It sets the partition, integration,
// and factory functions, as well as the state manager and options. The node listens for incoming signals,
// partitions the signal's data, processes the partitions, and integrates the results.
func NewSimplePartitionerNode(pfn node.PartitionerFn, ifn node.IntegratorFn, fac node.Factory, mgr node.StateManager, options node.Options) *SimplePartitionerNode {
	n := SimplePartitionerNode{
		partitionFunc:   pfn, // Set the partition function
		integrationFunc: ifn, // Set the integration function
		factory:         fac, // Set the factory function
	}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	// Start a goroutine to listen for signals and process them
	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.processSignal(sig) // Process the signal
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return // Terminate when done
			}
		}
	}()

	return &n
}

// SetPartitionFunc updates the partition function used by the node.
// This function defines how the signal's data will be divided into smaller parts.
func (n *SimplePartitionerNode) SetPartitionFunc(partitionFunc node.PartitionerFn) {
	n.partitionFunc = partitionFunc
}

// SetIntegrationFunc updates the integration function used by the node.
// This function defines how the results from the partitioned data will be combined into a final result.
func (n *SimplePartitionerNode) SetIntegrationFunc(integratonFunc node.IntegratorFn) {
	n.integrationFunc = integratonFunc
}

// SetNodeFactory updates the factory function used by the node.
// The factory is responsible for creating nodes that will process the partitioned data.
func (n *SimplePartitionerNode) SetNodeFactory(factory node.Factory) {
	n.factory = factory
}

// processSignal handles the signal processing for the SimplePartitionerNode. It first applies signal preprocessing,
// then partitions the signal's data, creates new nodes to process the partitions, and integrates the results.
// If any error occurs, the signal is marked as failed. Otherwise, the final integrated result is sent to connected nodes.
func (n *SimplePartitionerNode) processSignal(sig node.Signal) {
	var err error

	// Ensure that the required functions (partition, integration, and factory) are set
	if n.partitionFunc == nil || n.factory == nil || n.integrationFunc == nil {
		err := fmt.Errorf("partition, integrator, and factory functions are required")
		n.Fail(sig, err)
		return
	}

	// Preprocess the signal
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	sig.Status = StatusInProcess // Mark signal as in process

	// Partition the signal's task data into smaller parts
	parts, err := n.partitionFunc(sig.Task.String())
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}

	// Create a set of nodes to process the partitioned data
	nodes := n.factory(len(parts))
	respChan := make(chan node.Signal, len(parts)) // Channel to collect responses from the nodes
	emptyNode := &EmptyNode{inputCh: respChan}     // Empty node to gather results

	// Send each partitioned task to a separate node for processing
	for i, task := range parts {
		newSig := NewSignalFromSignal(nodes[i].ID(), sig)
		newSig.Task = &Carrier{TextData: task}
		nodes[i].Connect(emptyNode)  // Connect the node to the empty node
		nodes[i].InputCh() <- newSig // Send the signal to the node
	}

	// Collect the results from the nodes
	respList := []string{}
	for i := 0; i < len(parts); i++ {
		recSig := <-respChan
		respList = append(respList, recSig.Task.String()) // Collect the task results
	}

	// Integrate the results from the partitions
	response, err := n.integrationFunc(respList)
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}

	sig.Result = &Carrier{TextData: response}
	sig.Status = StatusSuccess

	// Post-process the signal
	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	// Send the processed signal to connected nodes with a 2-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure context is canceled after sending the signal

	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
}
