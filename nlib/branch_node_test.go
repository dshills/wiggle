package nlib_test

import (
	"testing"
	"time"

	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/nmock"
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/assert"
)

func TestSimpleBranchNode_AddConditional(t *testing.T) {
	mgr := new(nmock.MockStateManager)
	mgr.On("Register").Return(make((chan struct{})))
	options := node.Options{}
	branchNode := nlib.NewSimpleBranchNode(mgr, options)

	mockTargetNode := new(nmock.MockNode) // Use MockNode instead of MockGuidance
	mockCondFn := func(node.Signal) bool { return true }
	cond := node.BranchCondition{Target: mockTargetNode, ConditionFn: mockCondFn}
	branchNode.AddConditional(cond)

	// Ensure the condition was added correctly
	assert.Equal(t, 1, len(branchNode.Conditions()))
	assert.Equal(t, mockTargetNode, branchNode.Conditions()[0].Target)
	assert.NotNil(t, branchNode.Conditions()[0].ConditionFn)
}

func TestSimpleBranchNode_processSignal_ConditionalMatch(t *testing.T) {
	mgr := nlib.NewSimpleStateManager(nil)
	options := node.Options{}
	branchNode := nlib.NewSimpleBranchNode(mgr, options)

	mockTargetNode := new(nmock.MockNode)
	mockTargetNode.On("ID").Return("mockNode")
	mockInputCh := make(chan node.Signal, 1)
	mockTargetNode.On("InputCh").Return(mockInputCh)

	mockCondFn := func(node.Signal) bool { return true }
	cond := node.BranchCondition{Target: mockTargetNode, ConditionFn: mockCondFn}
	branchNode.AddConditional(cond)

	signal := node.Signal{NodeID: "test-signal"}

	go branchNode.ProcessSignal(signal)

	select {
	case receivedSignal := <-mockInputCh:
		assert.Equal(t, "mockNode", receivedSignal.NodeID)
	case <-time.After(2 * time.Second):
		t.Fatal("Signal was not sent to target node")
	}

	mockTargetNode.AssertCalled(t, "ID")
	mockTargetNode.AssertCalled(t, "InputCh")
}

func TestSimpleBranchNode_processSignal_NoMatch(t *testing.T) {
	mgr := nlib.NewSimpleStateManager(nil)
	options := node.Options{}
	branchNode := nlib.NewSimpleBranchNode(mgr, options)

	mockTargetNode := new(nmock.MockNode)
	mockInputCh := make(chan node.Signal, 1)
	mockTargetNode.On("InputCh").Return(mockInputCh)

	mockCondFn := func(node.Signal) bool { return false }
	cond := node.BranchCondition{Target: mockTargetNode, ConditionFn: mockCondFn}
	branchNode.AddConditional(cond)

	mockConnectedNode := new(nmock.MockNode)
	mockConnectedNode.On("InputCh").Return(mockInputCh)

	// Simulate SendToConnected
	mockConnectedNode.On("ID").Return("connectedNode")

	signal := node.Signal{NodeID: "test-signal"}
	branchNode.Connect(mockConnectedNode)

	go branchNode.ProcessSignal(signal)

	select {
	case receivedSignal := <-mockInputCh:
		assert.Equal(t, "connectedNode", receivedSignal.NodeID)
	case <-time.After(2 * time.Second):
		t.Fatal("Signal was not sent to connected node")
	}

	mockConnectedNode.AssertCalled(t, "ID")
	mockConnectedNode.AssertCalled(t, "InputCh")
}
