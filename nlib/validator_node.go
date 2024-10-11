package nlib

import (
	"context"
	"fmt"

	"github.com/dshills/wiggle/node"
	"github.com/dshills/wiggle/schema"
)

// Compile-time check
var _ node.Node = (*JSONValidatorNode)(nil)

type JSONValidatorNode struct {
	EmptyNode
	jsonSchema schema.Schema
}

func NewJSONValidatorNode(mgr node.StateManager, jsonSchema schema.Schema, options node.Options) *JSONValidatorNode {
	n := JSONValidatorNode{jsonSchema: jsonSchema}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.ProcessSignal(sig)
			case <-n.StateManager().Register():
				n.LogInfo("Received done")
			}
		}
	}()
	return &n
}

func (n *JSONValidatorNode) ProcessSignal(sig node.Signal) {
	var err error
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}
	sig.Status = StatusInProcess

	jmap, err := JSONAttemptMap(sig.Task.String())
	if err != nil {
		err = fmt.Errorf("JSONAttemptMap: %w", err)
		n.Fail(sig, err)
		return
	}

	if err := schema.Validate(jmap, n.jsonSchema, nil); err != nil {
		err = fmt.Errorf("schema.Validate: %w", err)
		n.Fail(sig, err)
		return
	}
	sig.Result = NewTextCarrier(sig.Task.String())

	sig.Status = StatusSuccess

	// Run post-processing hooks and forward the signal to connected nodes.
	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	if err := n.SendToConnected(context.TODO(), sig); err != nil {
		n.Fail(sig, err)
		return
	}
}
